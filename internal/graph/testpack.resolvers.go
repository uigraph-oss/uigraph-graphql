package graph

import (
	"context"
	"time"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) CreateTestPack(ctx context.Context, orgID string, serviceID string, input model.CreateTestPackInput) (*model.TestPack, error) {
	p, err := r.TestPack.CreateTestPack(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestPackToModel(p), nil
}

func (r *mutationResolver) UpdateTestPack(ctx context.Context, orgID string, serviceID string, id string, input model.UpdateTestPackInput) (*model.TestPack, error) {
	p, err := r.TestPack.UpdateTestPack(ctx, orgID, serviceID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestPackToModel(p), nil
}

func (r *mutationResolver) DeleteTestPack(ctx context.Context, orgID string, serviceID string, id string) (bool, error) {
	return true, r.TestPack.DeleteTestPack(ctx, orgID, serviceID, id)
}

func (r *mutationResolver) CreateTestCase(ctx context.Context, orgID string, serviceID string, input model.CreateTestCaseInput) (*model.TestCase, error) {
	tc, err := r.TestPack.CreateTestCase(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestCaseToModel(tc), nil
}

func (r *mutationResolver) UpdateTestCase(ctx context.Context, orgID string, serviceID string, id string, input model.UpdateTestCaseInput) (*model.TestCase, error) {
	tc, err := r.TestPack.UpdateTestCase(ctx, orgID, serviceID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestCaseToModel(tc), nil
}

func (r *mutationResolver) DeleteTestCase(ctx context.Context, orgID string, serviceID string, id string) (bool, error) {
	return true, r.TestPack.DeleteTestCase(ctx, orgID, serviceID, id)
}

func (r *mutationResolver) CreateTestRun(ctx context.Context, orgID string, serviceID string, input model.CreateTestRunInput) (*model.TestRun, error) {
	tr, err := r.TestPack.CreateTestRun(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestRunToModel(tr), nil
}

func (r *mutationResolver) UpdateTestRun(ctx context.Context, orgID string, serviceID string, id string, input model.UpdateTestRunInput) (*model.TestRun, error) {
	tr, err := r.TestPack.UpdateTestRun(ctx, orgID, serviceID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestRunToModel(tr), nil
}

func (r *mutationResolver) CreateTestRunResult(ctx context.Context, orgID string, serviceID string, input model.CreateTestRunResultInput) (*model.TestRunResult, error) {
	rr, err := r.TestPack.CreateTestRunResult(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestRunResultToModel(rr), nil
}

func (r *mutationResolver) UpdateTestRunResult(ctx context.Context, orgID string, serviceID string, id string, input model.UpdateTestRunResultInput) (*model.TestRunResult, error) {
	rr, err := r.TestPack.UpdateTestRunResult(ctx, orgID, serviceID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.TestRunResultToModel(rr), nil
}

func (r *queryResolver) TestPacks(ctx context.Context, orgID string, serviceID string) ([]*model.TestPack, error) {
	packs, err := r.TestPack.ListTestPacks(ctx, orgID, serviceID)
	if err != nil {
		return nil, err
	}
	return convert.TestPacksToModel(packs), nil
}

func (r *queryResolver) TestPackByID(ctx context.Context, orgID string, id string) (*model.TestPack, error) {
	p, err := r.TestPack.GetTestPackByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.TestPackToModel(p), nil
}

func (r *queryResolver) TestCases(ctx context.Context, orgID string, serviceID string, testPackID *string) ([]*model.TestCase, error) {
	cases, err := r.TestPack.ListTestCases(ctx, orgID, serviceID, testPackID)
	if err != nil {
		return nil, err
	}
	return convert.TestCasesToModel(cases), nil
}

func (r *queryResolver) TestRun(ctx context.Context, orgID string, serviceID string, id string) (*model.TestRun, error) {
	tr, err := r.TestPack.GetTestRun(ctx, orgID, serviceID, id)
	if err != nil {
		return nil, err
	}
	return convert.TestRunToModel(tr), nil
}

func (r *queryResolver) TestRuns(ctx context.Context, orgID string, serviceID string, testPackID *string) ([]*model.TestRun, error) {
	runs, err := r.TestPack.ListTestRuns(ctx, orgID, serviceID, testPackID)
	if err != nil {
		return nil, err
	}
	return convert.TestRunsToModel(runs), nil
}

func (r *queryResolver) TestRunsSummary(ctx context.Context, orgID string, serviceID string, testPackID *string, environment *string, status *string, executedBy *string, fromDate *time.Time, toDate *time.Time) ([]*model.TestRunSummary, error) {
	summary, err := r.TestPack.ListTestRunsSummary(ctx, orgID, serviceID, testPackID, environment, status, executedBy, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	return convert.TestRunSummariesToModel(summary), nil
}

func (r *queryResolver) TestRunResults(ctx context.Context, orgID string, serviceID string, testRunID string) ([]*model.TestRunResult, error) {
	results, err := r.TestPack.ListTestRunResults(ctx, orgID, serviceID, testRunID)
	if err != nil {
		return nil, err
	}
	return convert.TestRunResultsToModel(results), nil
}
