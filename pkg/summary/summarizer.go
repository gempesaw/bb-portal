package summary

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"

	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bescore"
)

// Summarizer struct.
type Summarizer struct {
	summary         *Summary
	problemDetector detectors.ProblemDetector
}

// Summarize function.
func Summarize(ctx context.Context, eventFileURL string) (*Summary, error) {
	reader, err := os.Open(eventFileURL)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", eventFileURL, err)
	}
	defer reader.Close()

	problemDetector := detectors.NewProblemDetector()
	summarizer := newSummarizer(eventFileURL, problemDetector)
	it := events.NewBuildEventIterator(ctx, reader)
	return summarizer.summarize(it)
}

// NewSummarizer constructor
func NewSummarizer() *Summarizer {
	return newSummarizer("", detectors.NewProblemDetector())
}

// newSummarizer
func newSummarizer(eventFileURL string, problemDetector detectors.ProblemDetector) *Summarizer {
	return &Summarizer{
		summary: &Summary{
			InvocationSummary: &InvocationSummary{},
			EventFileURL:      eventFileURL,
			RelatedFiles: map[string]string{
				filepath.Base(eventFileURL): eventFileURL,
			},
		},
		problemDetector: problemDetector,
	}
}

// summarize
func (s Summarizer) summarize(it *events.BuildEventIterator) (*Summary, error) {
	for {
		buildEvent, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get build event: %w", err)
		}

		err = s.ProcessEvent(buildEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to process event (with id: %s): %w", buildEvent.Id.String(), err)
		}
	}

	return s.FinishProcessing()
}

// FinishProcessing function
func (s Summarizer) FinishProcessing() (*Summary, error) {
	// If problems are ignored for the exit code, return immediately.
	if !shouldIgnoreProblems(s.summary.ExitCode) {
		// Add any detected test problems.
		problems, problemsErr := s.problemDetector.Problems()
		if problemsErr != nil {
			return nil, problemsErr
		}
		s.summary.Problems = append(s.summary.Problems, problems...)
	}

	return s.summary, nil
}

// ProcessEvent function
func (s Summarizer) ProcessEvent(buildEvent *events.BuildEvent) error {
	// Let problem detector process every event.
	s.problemDetector.ProcessBEPEvent(buildEvent)

	switch buildEvent.GetId().GetId().(type) {
	case *bes.BuildEventId_Started:
		s.handleStarted(buildEvent.GetStarted())

	case *bes.BuildEventId_BuildMetadata:
		s.handleBuildMetadata(buildEvent.GetBuildMetadata())

	case *bes.BuildEventId_BuildFinished:
		s.handleBuildFinished(buildEvent.GetFinished())

	case *bes.BuildEventId_BuildMetrics:
		s.handleBuildMetrics(buildEvent.GetBuildMetrics())

	case *bes.BuildEventId_StructuredCommandLine:
		err := s.handleStructuredCommandLine(buildEvent.GetStructuredCommandLine())
		if err != nil {
			return err
		}
	case *bes.BuildEventId_Configuration:
		s.handleBuildConfiguration(buildEvent.GetConfiguration())

	case *bes.BuildEventId_TargetConfigured:
		s.handleTargetConfigured(buildEvent.GetConfigured(), buildEvent.GetTargetConfiguredLabel(), time.Now())

	case *bes.BuildEventId_TargetCompleted:
		s.handleTargetCompleted(buildEvent.GetCompleted(), buildEvent.GetTargetCompletedLabel(), buildEvent.GetAborted(), time.Now())

	case *bes.BuildEventId_Fetch:
		s.handleFetch(buildEvent.GetFetch())

	case *bes.BuildEventId_TestResult:
		s.handleTestResult(buildEvent.GetTestResult(), buildEvent.GetId().GetTestResult().Label)

	case *bes.BuildEventId_TestSummary:
		s.handleTestSummary(buildEvent.GetTestSummary(), buildEvent.GetId().GetTestSummary().Label)

	case *bes.BuildEventId_OptionsParsed:
		s.handleOptionsParsed(buildEvent.GetOptionsParsed())

	case *bes.BuildEventId_BuildToolLogs:
		err := s.handleBuildToolLogs(buildEvent.GetBuildToolLogs())
		if err != nil {
			return err
		}
	case *bes.BuildEventId_Progress:
		s.handleProgress(buildEvent.GetProgress())
	}

	s.summary.BEPCompleted = buildEvent.GetLastMessage()
	return nil
}

// handleStarted
func (s Summarizer) handleStarted(started *bes.BuildStarted) {
	var startedAt time.Time
	if started.GetStartTime() != nil {
		startedAt = started.GetStartTime().AsTime()
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		startedAt = time.UnixMilli(started.GetStartTimeMillis())
	}
	s.summary.StartedAt = startedAt
	s.summary.InvocationID = started.GetUuid()
	s.summary.BazelVersion = started.GetBuildToolVersion()
}

// handleFetch
func (s Summarizer) handleFetch(fetch *bes.Fetch) {
	if fetch.Success {
		s.summary.NumFetches++
	}
}

// handleBuildConfiguration
func (s Summarizer) handleBuildConfiguration(configuration *bes.Configuration) {
	s.summary.CPU = configuration.Cpu
	s.summary.PlatformName = configuration.PlatformName
	s.summary.ConfigrationMnemonic = configuration.Mnemonic
}

// handleTargetConfigured
func (s Summarizer) handleTargetConfigured(target *bes.TargetConfigured, label string, timestamp time.Time) {
	if len(label) == 0 {
		panic("missing a target label for target configured event!")
	}
	if target == nil {
		slog.Debug(fmt.Sprintf("missing target for label %s on targetConfigured", label))
		return
	}

	// if this is the first target we've seen, initialize the targets collection
	if s.summary.Targets == nil {
		s.summary.Targets = make(map[string]TargetPair)
	}

	// create a target pair and at it to the targets collection
	s.summary.Targets[label] = TargetPair{
		Configuration: TargetConfigured{
			StartTimeInMs: timestamp.UnixMilli(),
			TargetKind:    target.TargetKind,
			TestSize:      TestSize(target.TestSize),
			Tag:           target.Tag,
		},
		Success:    false, // set it to false, change it when we get a complete
		TargetKind: target.TargetKind,
		TestSize:   TestSize(target.TestSize),
	}
}

// handleTargetCompleted
func (s Summarizer) handleTargetCompleted(target *bes.TargetComplete, label string, aborted *bes.Aborted, timestamp time.Time) {
	if len(label) == 0 {
		panic("label is empty for a target completed event")
	}

	if s.summary.Targets == nil {
		panic(fmt.Sprintf("target completed event received before any target configured messages for label %s,", label))
	}

	var targetPair TargetPair
	targetPair, ok := s.summary.Targets[label]

	if !ok {
		panic(fmt.Sprintf("target completed event received for label %s before target configured message received", label))
	}

	var targetCompletion TargetComplete

	if target != nil {
		targetCompletion = TargetComplete{
			Success:     target.Success,
			Tag:         target.Tag,
			EndTimeInMs: timestamp.UnixMilli(),
		}
		if target.TestTimeout != nil {
			targetCompletion.TestTimeoutSeconds = target.TestTimeout.Seconds
			targetCompletion.TestTimeout = target.TestTimeout.Seconds
		}
	} else { // this event was aborted
		targetCompletion = TargetComplete{
			Success:     false,
			EndTimeInMs: timestamp.UnixMilli(),
		}
	}

	targetPair.Completion = targetCompletion
	targetPair.DurationInMs = targetPair.Completion.EndTimeInMs - targetPair.Configuration.StartTimeInMs
	targetPair.Success = targetCompletion.Success

	if aborted != nil {
		targetPair.AbortReason = AbortReason(aborted.Reason)
	}

	s.summary.Targets[label] = targetPair
}

// handleTestResult
func (s Summarizer) handleTestResult(testResult *bes.TestResult, label string) {
	if len(label) == 0 {
		panic("missing label on TestResult event")
	}
	if testResult == nil {
		panic(fmt.Sprintf("missing TestResult for label %s", label))
	}
	var testResults []TestResult
	if s.summary.Tests == nil {
		s.summary.Tests = make(map[string]TestsCollection)
	}
	testcollection, ok := s.summary.Tests[label]
	if ok {
		testResults = testcollection.TestResults
	} else { // initailize it if we've never seen this label before
		testcollection = TestsCollection{
			TestSummary:    TestSummary{},
			TestResults:    []TestResult{},
			CachedLocally:  true,
			CachedRemotely: true,
			Strategy:       "INITIALIZED",
		}
		testResults = make([]TestResult, 0)
	}
	tr := TestResult{
		Status:              TestStatus(testResult.Status),
		StatusDetails:       testResult.StatusDetails,
		Label:               label,
		Warning:             testResult.Warning,
		CachedLocally:       testResult.CachedLocally,
		TestAttemptStart:    testResult.TestAttemptStart.AsTime().String(),
		TestAttemptDuration: testResult.TestAttemptDuration.AsDuration().Milliseconds(),
		ExecutionInfo:       processExecutionInfo(testResult),
		TestActionOutput:    make([]TestFile, 0),
	}
	for _, ao := range testResult.TestActionOutput {
		actionOutput := TestFile{
			Digest: ao.Digest,
			File:   ao.GetUri(),
			Length: ao.Length,
			Name:   ao.Name,
			Prefix: ao.PathPrefix,
		}
		tr.TestActionOutput = append(tr.TestActionOutput, actionOutput)
	}
	testResults = append(testResults, tr)
	testcollection.TestResults = testResults
	if testResult.Status != bes.TestStatus_NO_STATUS {
		if !testResult.CachedLocally {
			testcollection.CachedLocally = false
		}
		if !tr.ExecutionInfo.CachedRemotely {
			testcollection.CachedRemotely = false
		}
		if (testcollection.Strategy) == "INITIALIZED" {
			testcollection.Strategy = tr.ExecutionInfo.Strategy
		} else {
			if testcollection.Strategy != tr.ExecutionInfo.Strategy {
				testcollection.Strategy = "indeterminate"
			}
		}
	}

	s.summary.Tests[label] = testcollection
}

// processExecutionInfo
func processExecutionInfo(testResult *bes.TestResult) ExecutionInfo {
	var result ExecutionInfo
	var timingBreakdown TimingBreakdown
	var children []TimingChild
	if testResult.ExecutionInfo != nil {
		if testResult.ExecutionInfo.TimingBreakdown != nil {
			for _, c := range testResult.ExecutionInfo.TimingBreakdown.Child {
				child := TimingChild{
					Name: c.Name,
					Time: c.Time.AsDuration().String(),
				}
				children = append(children, child)
			}

			timingBreakdown.Name = testResult.ExecutionInfo.TimingBreakdown.Name
			timingBreakdown.Time = testResult.ExecutionInfo.TimingBreakdown.Time.String()
			timingBreakdown.Child = children
		}

		result.Strategy = testResult.ExecutionInfo.Strategy
		result.CachedRemotely = testResult.ExecutionInfo.CachedRemotely
		result.ExitCode = testResult.ExecutionInfo.ExitCode
		result.Hostname = testResult.ExecutionInfo.Hostname
		result.TimingBreakdown = timingBreakdown
	}
	return result
}

// handleTestSummary
func (s Summarizer) handleTestSummary(testSummary *bes.TestSummary, label string) {
	if len(label) == 0 {
		panic("missing label on handleTestSummary event")
	}

	if testSummary == nil {
		panic(fmt.Sprintf("missing test summary object for handleTestSummary event for label %s", label))
	}

	testCollection, ok := s.summary.Tests[label]

	if !ok {
		panic(fmt.Sprintf("received a test summary event but never first saw a test result for label %s", label))
	}

	tSummary := testCollection.TestSummary

	tSummary.AttemptCount = testSummary.AttemptCount
	tSummary.FirstStartTime = testSummary.FirstStartTime.AsTime().Unix()
	tSummary.Label = label
	tSummary.LastStopTime = testSummary.FirstStartTime.AsTime().Unix()
	tSummary.RunCount = testSummary.RunCount
	tSummary.ShardCount = testSummary.ShardCount
	tSummary.Status = TestStatus(testSummary.OverallStatus)
	tSummary.TotalNumCached = testSummary.TotalNumCached
	tSummary.TotalRunCount = testSummary.TotalRunCount
	tSummary.TotalRunDuration = testSummary.TotalRunDuration.AsDuration().Microseconds()

	testCollection.TestSummary = tSummary
	testCollection.OverallStatus = tSummary.Status
	testCollection.DurationMs = tSummary.TotalRunDuration
	s.summary.Tests[label] = testCollection
}

// handleBuildMetadata
func (s Summarizer) handleBuildMetadata(metadataProto *bes.BuildMetadata) {
	metadataMap := metadataProto.GetMetadata()
	// extract user data
	if metadataMap == nil {
		return
	}
	stepLabel, stepLabelOk := metadataMap[stepLabelKey]
	if !stepLabelOk {
		slog.Debug("No step label found in build metadata")
	}
	userEmail, userEmailOk := metadataMap[userEmailKey]
	if !userEmailOk {
		slog.Debug("No user email found in build metadata")
	}
	userLdap, userLdapOk := metadataMap[userLdapKey]
	if !userLdapOk {
		slog.Debug("No user ldap information found in build metadata")
	}
	s.summary.StepLabel = stepLabel
	s.summary.UserEmail = userEmail
	s.summary.UserLDAP = userLdap
}

// handleBuildMetrics
func (s Summarizer) handleBuildMetrics(metrics *bes.BuildMetrics) {
	actionSummary := readActionSummary(metrics.ActionSummary)
	memoryMetrics := readMemoryMetrics(metrics.MemoryMetrics)
	targetMetrics := readTargetMetrics(metrics.TargetMetrics)
	packageMetrics := readPackageMetrics(metrics.PackageMetrics)
	timingMetrics := readTimingMetrics(metrics.TimingMetrics)
	artifactMetrics := readArtifactMetrics(metrics.ArtifactMetrics)
	cumulativeMetrics := readCumulativeMetrics(metrics.CumulativeMetrics)
	networkMetrics := readNetworkMetrics(metrics.NetworkMetrics)
	buildGraphMetrics := readBuildGraphMetrics(metrics.BuildGraphMetrics)

	summaryMetrics := Metrics{
		ActionSummary:     actionSummary,
		MemoryMetrics:     memoryMetrics,
		TargetMetrics:     targetMetrics,
		PackageMetrics:    packageMetrics,
		TimingMetrics:     timingMetrics,
		ArtifactMetrics:   artifactMetrics,
		CumulativeMetrics: cumulativeMetrics,
		NetworkMetrics:    networkMetrics,
		BuildGraphMetrics: buildGraphMetrics,
		// DynamicExecutionMetrics: dynamicMetrics,
	}

	s.summary.Metrics = summaryMetrics
}

// readBuildGraphMetrics
func readBuildGraphMetrics(buildGraphMetricsData *bes.BuildMetrics_BuildGraphMetrics) BuildGraphMetrics {
	// TODO: these values are not on the proto currently.  once they are, update this code to pull them out
	// var dirtiedValues = make([]EvaluationStat, 0)
	// var changedValues = make([]EvaluationStat, 0)
	// var builtValues = make([]EvaluationStat, 0)
	// var cleanedValues = make([]EvaluationStat, 0)
	// var evaluatedValues = make([]EvaluationStat, 0)

	// DirtiedValues:                   dirtiedValues,
	// ChangedValues:                   changedValues,
	// BuiltValues:                     builtValues,
	// CleanedValues:                   cleanedValues,
	// EvaluatedValues:                 evaluatedValues,
	buildGraphMetrics := BuildGraphMetrics{
		ActionLookupValueCount:                    buildGraphMetricsData.ActionLookupValueCount,
		ActionLookupValueCountNotIncludingAspects: buildGraphMetricsData.ActionLookupValueCountNotIncludingAspects,
		ActionCount:                     buildGraphMetricsData.ActionCount,
		InputFileConfiguredTargetCount:  buildGraphMetricsData.InputFileConfiguredTargetCount,
		OutputFileConfiguredTargetCount: buildGraphMetricsData.OutputFileConfiguredTargetCount,
		OtherConfiguredTargetCount:      buildGraphMetricsData.OtherConfiguredTargetCount,
		OutputArtifactCount:             buildGraphMetricsData.OutputArtifactCount,
		PostInvocationSkyframeNodeCount: buildGraphMetricsData.PostInvocationSkyframeNodeCount,
	}
	return buildGraphMetrics
}

// readNetworkMetrics
func readNetworkMetrics(networkMetricsData *bes.BuildMetrics_NetworkMetrics) NetworkMetrics {
	if networkMetricsData == nil {
		return NetworkMetrics{}
	}
	systemNetworkStats := readSystemNetworkStats(networkMetricsData.SystemNetworkStats)

	networkMetrics := NetworkMetrics{
		SystemNetworkStats: &systemNetworkStats,
	}
	return networkMetrics
}

// readSystemNetworkStats
func readSystemNetworkStats(systemNetworkStatsData *bes.BuildMetrics_NetworkMetrics_SystemNetworkStats) SystemNetworkStats {
	var systemNetworkStats SystemNetworkStats
	if systemNetworkStatsData != nil {
		systemNetworkStats = SystemNetworkStats{
			BytesSent:             systemNetworkStatsData.BytesSent,
			BytesRecv:             systemNetworkStatsData.BytesRecv,
			PacketsSent:           systemNetworkStatsData.PacketsSent,
			PacketsRecv:           systemNetworkStatsData.PacketsRecv,
			PeakBytesSentPerSec:   systemNetworkStatsData.PeakBytesSentPerSec,
			PeakBytesRecvPerSec:   systemNetworkStatsData.PeakBytesRecvPerSec,
			PeakPacketsSentPerSec: systemNetworkStatsData.PeakPacketsSentPerSec,
			PeakPacketsRecvPerSec: systemNetworkStatsData.PeakPacketsRecvPerSec,
		}
	}
	return systemNetworkStats
}

// readCumulativeMetrics
func readCumulativeMetrics(cumulativeMetricsData *bes.BuildMetrics_CumulativeMetrics) CumulativeMetrics {
	cumulativeMetrics := CumulativeMetrics{
		NumAnalyses: cumulativeMetricsData.NumAnalyses,
		NumBuilds:   cumulativeMetricsData.NumBuilds,
	}
	return cumulativeMetrics
}

// readArtifactMetrics
func readArtifactMetrics(artifactMetricsData *bes.BuildMetrics_ArtifactMetrics) ArtifactMetrics {
	sourceArtifactsRead := FilesMetric{
		SizeInBytes: artifactMetricsData.SourceArtifactsRead.SizeInBytes,
		Count:       artifactMetricsData.SourceArtifactsRead.Count,
	}

	outputArtifactsSeen := FilesMetric{
		SizeInBytes: artifactMetricsData.OutputArtifactsSeen.SizeInBytes,
		Count:       artifactMetricsData.OutputArtifactsSeen.Count,
	}

	outputArtifactsFromActionCache := FilesMetric{
		SizeInBytes: artifactMetricsData.OutputArtifactsFromActionCache.SizeInBytes,
		Count:       artifactMetricsData.OutputArtifactsFromActionCache.Count,
	}

	topLevelArtifacts := FilesMetric{
		SizeInBytes: artifactMetricsData.TopLevelArtifacts.SizeInBytes,
		Count:       artifactMetricsData.TopLevelArtifacts.Count,
	}

	artifactMetrics := ArtifactMetrics{
		SourceArtifactsRead:            sourceArtifactsRead,
		OutputArtifactsSeen:            outputArtifactsSeen,
		OutputArtifactsFromActionCache: outputArtifactsFromActionCache,
		TopLevelArtifacts:              topLevelArtifacts,
	}
	return artifactMetrics
}

// readTimingMetrics
func readTimingMetrics(timingMetricsData *bes.BuildMetrics_TimingMetrics) TimingMetrics {
	timingMetrics := TimingMetrics{
		CPUTimeInMs:            timingMetricsData.CpuTimeInMs,
		WallTimeInMs:           timingMetricsData.WallTimeInMs,
		ExecutionPhaseTimeInMs: timingMetricsData.ExecutionPhaseTimeInMs,
		AnalysisPhaseTimeInMs:  timingMetricsData.AnalysisPhaseTimeInMs,
	}
	return timingMetrics
}

// readPackageMetrics
func readPackageMetrics(packagMetricsData *bes.BuildMetrics_PackageMetrics) PackageMetrics {
	packageLoadMetrics := readPackageLoadMetrics(packagMetricsData.PackageLoadMetrics)

	packageMetrics := PackageMetrics{
		PackagesLoaded:     packagMetricsData.PackagesLoaded,
		PackageLoadMetrics: packageLoadMetrics,
	}
	return packageMetrics
}

// readTargetMetrics
func readTargetMetrics(targetMetricsData *bes.BuildMetrics_TargetMetrics) TargetMetrics {
	targetMetrics := TargetMetrics{
		TargetsConfigured:                    targetMetricsData.TargetsConfigured,
		TargetsConfiguredNotIncludingAspects: targetMetricsData.TargetsConfiguredNotIncludingAspects,
		TargetsLoaded:                        targetMetricsData.TargetsLoaded,
	}
	return targetMetrics
}

// readMemoryMetrics
func readMemoryMetrics(memoryMetricsData *bes.BuildMetrics_MemoryMetrics) MemoryMetrics {
	garbageMetrics := readGarbageMetrics(memoryMetricsData.GarbageMetrics)

	memoryMetrics := MemoryMetrics{
		PeakPostGcHeapSize:             memoryMetricsData.PeakPostGcHeapSize,
		PeakPostGcTenuredSpaceHeapSize: memoryMetricsData.PeakPostGcTenuredSpaceHeapSize,
		UsedHeapSizePostBuild:          memoryMetricsData.UsedHeapSizePostBuild,
		GarbageMetrics:                 garbageMetrics,
	}
	return memoryMetrics
}

// readActionSummary
func readActionSummary(actionSummaryData *bes.BuildMetrics_ActionSummary) ActionSummary {
	actionCacheStatistics := readActionCacheStatistics(actionSummaryData.ActionCacheStatistics)

	runnerCounts := readRunnerCounts(actionSummaryData.RunnerCount)

	actionDatas := readActionDatas(actionSummaryData.ActionData)

	actionSummary := ActionSummary{
		ActionsCreated:                    actionSummaryData.ActionsCreated,
		ActionsExecuted:                   actionSummaryData.ActionsExecuted,
		ActionsCreatedNotIncludingAspects: actionSummaryData.ActionsCreatedNotIncludingAspects,
		ActionCacheStatistics:             actionCacheStatistics,
		RunnerCount:                       runnerCounts,
		ActionData:                        actionDatas,
	}
	return actionSummary
}

// readActionCacheStatistics
func readActionCacheStatistics(actionCacheStatisticsData *bescore.ActionCacheStatistics) ActionCacheStatistics {
	missDetails := readMissDetails(actionCacheStatisticsData.MissDetails)
	actionCacheStatistics := ActionCacheStatistics{
		SizeInBytes:  actionCacheStatisticsData.SizeInBytes,
		SaveTimeInMs: actionCacheStatisticsData.SaveTimeInMs,

		Hits:        actionCacheStatisticsData.Hits,
		Misses:      actionCacheStatisticsData.Misses,
		MissDetails: missDetails,
	}
	return actionCacheStatistics
}

// readPackageLoadMetrics
func readPackageLoadMetrics(packageLoadMetricsData []*bescore.PackageLoadMetrics) []PackageLoadMetrics {
	packageLoadMetrics := make([]PackageLoadMetrics, len(packageLoadMetricsData))

	for i, plm := range packageLoadMetricsData {
		packageLoadMetric := PackageLoadMetrics{
			Name:               *plm.Name,
			NumTargets:         *plm.NumTargets,
			LoadDuration:       plm.LoadDuration.AsDuration().Milliseconds(),
			ComputationSteps:   *plm.ComputationSteps,
			NumTransitiveLoads: *plm.NumTransitiveLoads,
			PackageOverhead:    *plm.PackageOverhead,
		}
		packageLoadMetrics[i] = packageLoadMetric
	}
	return packageLoadMetrics
}

// readGarbageMetrics
func readGarbageMetrics(garbageMetricsData []*bes.BuildMetrics_MemoryMetrics_GarbageMetrics) []GarbageMetrics {
	garbageMetrics := make([]GarbageMetrics, len(garbageMetricsData))

	for i, gm := range garbageMetricsData {
		garbageMetric := GarbageMetrics{
			Type:             gm.Type,
			GarbageCollected: gm.GarbageCollected,
		}
		garbageMetrics[i] = garbageMetric
	}
	return garbageMetrics
}

// readActionDatas
func readActionDatas(actionDataData []*bes.BuildMetrics_ActionSummary_ActionData) []ActionData {
	actionDatas := make([]ActionData, len(actionDataData))
	for i, ad := range actionDataData {
		actionData := ActionData{
			Mnemonic:        ad.Mnemonic,
			UserTime:        ad.UserTime.AsDuration().Milliseconds(),
			SystemTime:      ad.SystemTime.AsDuration().Milliseconds(),
			ActionsExecuted: ad.ActionsExecuted,
			FirstStartedMs:  ad.FirstStartedMs,
			LastEndedMs:     ad.LastEndedMs,
		}
		actionDatas[i] = actionData
	}
	return actionDatas
}

// readRunnerCounts
func readRunnerCounts(runnerCountsData []*bes.BuildMetrics_ActionSummary_RunnerCount) []RunnerCount {
	runnerCounts := make([]RunnerCount, len(runnerCountsData))
	for i, rc := range runnerCountsData {
		runnerCount := RunnerCount{
			ExecKind: rc.ExecKind,
			Count:    rc.Count,
			Name:     rc.Name,
		}
		runnerCounts[i] = runnerCount
	}
	return runnerCounts
}

// readMissDetails
func readMissDetails(missDetailsData []*bescore.ActionCacheStatistics_MissDetail) []MissDetail {
	missDetails := make([]MissDetail, len(missDetailsData))
	for i, md := range missDetailsData {
		missDetail := MissDetail{
			Count:  md.Count,
			Reason: MissReason(*md.Reason.Enum()),
		}
		missDetails[i] = missDetail
	}
	return missDetails
}

// handleBuildFinished
func (s Summarizer) handleBuildFinished(finished *bes.BuildFinished) {
	var endedAt time.Time
	if finished.GetFinishTime() != nil {
		endedAt = finished.GetFinishTime().AsTime()
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		endedAt = time.UnixMilli(finished.GetFinishTimeMillis())
	}
	s.summary.EndedAt = &endedAt
	s.summary.InvocationSummary.ExitCode = &ExitCode{
		Code: int(finished.GetExitCode().GetCode()),
		Name: finished.GetExitCode().GetName(),
	}
}

// handleStructuredCommandLine
func (s Summarizer) handleStructuredCommandLine(structuredCommandLine *bescore.CommandLine) error {
	if structuredCommandLine.GetCommandLineLabel() != "original" {
		return nil
	}

	s.updateEnvVarsAndCommandFromStructuredCommandLine(structuredCommandLine)

	// Parse Gerrit change number if available.
	if changeNumberStr, ok := s.summary.InvocationSummary.EnvVars["GERRIT_CHANGE_NUMBER"]; ok && changeNumberStr != "" {
		changeNumber, err := envToI(s.summary.InvocationSummary.EnvVars, "GERRIT_CHANGE_NUMBER")
		if err != nil {
			return err
		}
		s.summary.ChangeNumber = changeNumber
	}

	// Parse Gerrit patchset number if available.
	if patchsetNumberStr, ok := s.summary.InvocationSummary.EnvVars["GERRIT_PATCHSET_NUMBER"]; ok && patchsetNumberStr != "" {
		patchsetNumber, err := envToI(s.summary.InvocationSummary.EnvVars, "GERRIT_PATCHSET_NUMBER")
		if err != nil {
			return err
		}
		s.summary.PatchsetNumber = patchsetNumber
	}

	// Decode commit message, so that client doesn't have to.
	commitMessage := s.summary.InvocationSummary.EnvVars["GERRIT_CHANGE_COMMIT_MESSAGE"]
	if commitMessage != "" {
		decodedCommitMessage, err := base64.StdEncoding.DecodeString(commitMessage)
		if err == nil {
			s.summary.InvocationSummary.EnvVars["GERRIT_CHANGE_COMMIT_MESSAGE"] = string(decodedCommitMessage)
		} else {
			slog.Debug("GERRIT_CHANGE_COMMIT_MESSAGE was not base64 encoded, assuming it is normal string")
		}
	}

	// Set build URL and UUID
	s.summary.BuildURL = s.summary.InvocationSummary.EnvVars["BUILD_URL"]
	s.summary.BuildUUID = uuid.NewSHA1(uuid.NameSpaceURL, []byte(s.summary.BuildURL))

	return nil
}

// handleOptionsParsed
func (s Summarizer) handleOptionsParsed(optionsParsed *bes.OptionsParsed) {
	s.summary.InvocationSummary.BazelCommandLine.Options = optionsParsed.GetExplicitCmdLine()
}

// handleProgress
func (s Summarizer) handleProgress(progressMsg *bes.Progress) {
	s.summary.BuildLogs.WriteString(progressMsg.GetStderr())
	s.summary.BuildLogs.WriteString(progressMsg.GetStdout())
}

// handleBuildToolLogs
func (s Summarizer) handleBuildToolLogs(buildToolLogs *bes.BuildToolLogs) error {
	for _, logs := range buildToolLogs.GetLog() {
		uri := logs.GetUri()
		blobURI := detectors.BlobURI(uri)

		if s.summary.RelatedFiles == nil {
			s.summary.RelatedFiles = map[string]string{}
		}
		if logs.GetUri() != "" {
			s.summary.RelatedFiles[logs.GetName()] = string(blobURI)
		}
	}
	return nil
}

// updateEnvVarsAndCommandFromStructuredCommandLine
func (s Summarizer) updateEnvVarsAndCommandFromStructuredCommandLine(structuredCommandLine *bescore.CommandLine) {
	sections := structuredCommandLine.GetSections()
	for _, section := range sections {
		label := section.GetSectionLabel()
		if label == "command options" {
			s.summary.InvocationSummary.EnvVars = map[string]string{}
			ParseEnvVarsFromSectionOptions(section, &s.summary.InvocationSummary.EnvVars)
		} else if section.GetChunkList() != nil {
			sectionChunksStr := strings.Join(section.GetChunkList().GetChunk(), " ")
			switch label {
			case "executable":
				s.summary.InvocationSummary.BazelCommandLine.Executable = sectionChunksStr
			case "command":
				s.summary.InvocationSummary.BazelCommandLine.Command = sectionChunksStr
			case "residual":
				s.summary.InvocationSummary.BazelCommandLine.Residual = sectionChunksStr
			}
		}
	}
}

// shouldIgnoreProblems
func shouldIgnoreProblems(exitCode *ExitCode) bool {
	return exitCode != nil && (exitCode.Code == ExitCodeSuccess || exitCode.Code == ExitCodeInterrupted)
}

// envToI
func envToI(envVars map[string]string, name string) (int, error) {
	res, err := strconv.Atoi(envVars[name])
	if err != nil {
		slog.Error("failed to parse env var to int", "envKey", name, "envValue", envVars[name], "err", err)
		return 0, fmt.Errorf("failed to parse %s (value: %s) as an int: %w", name, envVars[name], err)
	}
	return res, nil
}

// ParseEnvVarsFromSectionOptions function
func ParseEnvVarsFromSectionOptions(section *bescore.CommandLineSection, destMap *map[string]string) {
	if section.GetOptionList() == nil {
		return
	}
	options := section.GetOptionList().GetOption()
	for _, option := range options {
		if option.GetOptionName() != "client_env" {
			// Only looking for env vars from the client env
			continue
		}
		envPair := option.GetOptionValue()
		equalIndex := strings.Index(envPair, "=")
		if equalIndex <= 0 {
			// Skip anything missing an equals sign. The env vars come in the format key=value
			continue
		}
		envName := envPair[:equalIndex]
		envValue := envPair[equalIndex+1:]
		(*destMap)[envName] = envValue
	}
}
