package ghbot

import (
	"context"

	"github.com/google/go-github/v25/github"
)

type CheckRunEventHook func(context.Context, *github.CheckRunEvent) error
type CheckSuiteEventHook func(context.Context, *github.CheckSuiteEvent) error
type CommitCommentEventHook func(context.Context, *github.CommitCommentEvent) error
type CreateEventHook func(context.Context, *github.CreateEvent) error
type DeleteEventHook func(context.Context, *github.DeleteEvent) error
type DeployKeyEventHook func(context.Context, *github.DeployKeyEvent) error
type DeploymentEventHook func(context.Context, *github.DeploymentEvent) error
type DeploymentStatusEventHook func(context.Context, *github.DeploymentStatusEvent) error
type ForkEventHook func(context.Context, *github.ForkEvent) error
type GitHubAppAuthorizationEventHook func(context.Context, *github.GitHubAppAuthorizationEvent) error
type GollumEventHook func(context.Context, *github.GollumEvent) error
type InstallationEventHook func(context.Context, *github.InstallationEvent) error
type InstallationRepositoriesEventHook func(context.Context, *github.InstallationRepositoriesEvent) error
type IssueCommentEventHook func(context.Context, *github.IssueCommentEvent) error
type IssueEventHook func(context.Context, *github.IssueEvent) error
type IssuesEventHook func(context.Context, *github.IssuesEvent) error
type LabelEventHook func(context.Context, *github.LabelEvent) error
type MarketplacePurchaseEventHook func(context.Context, *github.MarketplacePurchaseEvent) error
type MemberEventHook func(context.Context, *github.MemberEvent) error
type MembershipEventHook func(context.Context, *github.MembershipEvent) error
type MetaEventHook func(context.Context, *github.MetaEvent) error
type MilestoneEventHook func(context.Context, *github.MilestoneEvent) error
type OrgBlockEventHook func(context.Context, *github.OrgBlockEvent) error
type OrganizationEventHook func(context.Context, *github.OrganizationEvent) error
type PageBuildEventHook func(context.Context, *github.PageBuildEvent) error
type PingEventHook func(context.Context, *github.PingEvent) error
type ProjectCardEventHook func(context.Context, *github.ProjectCardEvent) error
type ProjectColumnEventHook func(context.Context, *github.ProjectColumnEvent) error
type ProjectEventHook func(context.Context, *github.ProjectEvent) error
type PublicEventHook func(context.Context, *github.PublicEvent) error
type PullRequestEventHook func(context.Context, *github.PullRequestEvent) error
type PullRequestReviewCommentEventHook func(context.Context, *github.PullRequestReviewCommentEvent) error
type PullRequestReviewEventHook func(context.Context, *github.PullRequestReviewEvent) error
type PushEventHook func(context.Context, *github.PushEvent) error
type ReleaseEventHook func(context.Context, *github.ReleaseEvent) error
type RepositoryEventHook func(context.Context, *github.RepositoryEvent) error
type RepositoryVulnerabilityAlertEventHook func(context.Context, *github.RepositoryVulnerabilityAlertEvent) error
type StarEventHook func(context.Context, *github.StarEvent) error
type StatusEventHook func(context.Context, *github.StatusEvent) error
type TeamAddEventHook func(context.Context, *github.TeamAddEvent) error
type TeamEventHook func(context.Context, *github.TeamEvent) error
type WatchEventHook func(context.Context, *github.WatchEvent) error
