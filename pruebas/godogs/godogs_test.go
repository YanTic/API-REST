package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
)

type godogsCtxKey struct{}

func iEat(ctx context.Context, num int) (context.Context, error) {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return ctx, errors.New("There are no godogs available")
	}

	if available < num {
		return ctx, fmt.Errorf("you cannot eat %d godogs, there are %d available", num, available)
	}

	available -= num

	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func thereAreGodogs(ctx context.Context, available int) (context.Context, error) {
	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return errors.New("There are no godogs available")
	}

	if available != remaining {
		return fmt.Errorf("you cannot eat %d godogs, there are %d available", remaining, available)
	}

	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^there are (\d+) godogs$`, thereAreGodogs)
	ctx.When(`^I eat (\d+)$`, iEat)
	ctx.Then(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}

// func InitializeScenario(ctx *godog.ScenarioContext) {
// 	ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
// 	ctx.Step(`^I eat (\d+)$`, iEat)
// 	ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
// }
