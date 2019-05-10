package ghbot

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-github/v25/github"
	"golang.org/x/xerrors"
)

type Bot struct {
	mu            sync.Mutex
	webhookSecret []byte
	logger        Logger

	checkRunEventHooks                     []CheckRunEventHook
	checkSuiteEventHooks                   []CheckSuiteEventHook
	commitCommentEventHooks                []CommitCommentEventHook
	createEventHooks                       []CreateEventHook
	deleteEventHooks                       []DeleteEventHook
	deployKeyEventHooks                    []DeployKeyEventHook
	deploymentEventHooks                   []DeploymentEventHook
	deploymentStatusEventHooks             []DeploymentStatusEventHook
	forkEventHooks                         []ForkEventHook
	gitHubAppAuthorizationEventHooks       []GitHubAppAuthorizationEventHook
	gollumEventHooks                       []GollumEventHook
	installationEventHooks                 []InstallationEventHook
	installationRepositoriesEventHooks     []InstallationRepositoriesEventHook
	issueCommentEventHooks                 []IssueCommentEventHook
	issueEventHooks                        []IssueEventHook
	issuesEventHooks                       []IssuesEventHook
	labelEventHooks                        []LabelEventHook
	marketplacePurchaseEventHooks          []MarketplacePurchaseEventHook
	memberEventHooks                       []MemberEventHook
	membershipEventHooks                   []MembershipEventHook
	metaEventHooks                         []MetaEventHook
	milestoneEventHooks                    []MilestoneEventHook
	orgBlockEventHooks                     []OrgBlockEventHook
	organizationEventHooks                 []OrganizationEventHook
	pageBuildEventHooks                    []PageBuildEventHook
	pingEventHooks                         []PingEventHook
	projectCardEventHooks                  []ProjectCardEventHook
	projectColumnEventHooks                []ProjectColumnEventHook
	projectEventHooks                      []ProjectEventHook
	publicEventHooks                       []PublicEventHook
	pullRequestEventHooks                  []PullRequestEventHook
	pullRequestReviewCommentEventHooks     []PullRequestReviewCommentEventHook
	pullRequestReviewEventHooks            []PullRequestReviewEventHook
	pushEventHooks                         []PushEventHook
	releaseEventHooks                      []ReleaseEventHook
	repositoryEventHooks                   []RepositoryEventHook
	repositoryVulnerabilityAlertEventHooks []RepositoryVulnerabilityAlertEventHook
	starEventHooks                         []StarEventHook
	statusEventHooks                       []StatusEventHook
	teamAddEventHooks                      []TeamAddEventHook
	teamEventHooks                         []TeamEventHook
	watchEventHooks                        []WatchEventHook
}

type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

type nopLogger struct{}

func (l nopLogger) Print(...interface{})          {}
func (l nopLogger) Printf(string, ...interface{}) {}
func (l nopLogger) Println(...interface{})        {}

type Config struct {
	WebHookSecret string
}

func New(cfg Config) *Bot {
	bot := Bot{
		webhookSecret: []byte(cfg.WebHookSecret),
		logger:        nopLogger{},
	}
	return &bot
}

func (bot *Bot) SetLogger(logger Logger) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	bot.logger = logger
}

func (bot *Bot) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/github", bot.githubWebHookHandler)
	httpSrv := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	return httpSrv.ListenAndServe()
}

func (bot *Bot) githubWebHookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	payload, err := github.ValidatePayload(r, bot.webhookSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := bot.handleWebHookEvent(r.Context(), github.WebHookType(r), event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (bot *Bot) handleWebHookEvent(ctx context.Context, typ string, event interface{}) error {
	switch e := event.(type) {
	case *github.CheckRunEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.checkRunEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.CheckSuiteEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.checkSuiteEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.CommitCommentEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.commitCommentEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.CreateEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.createEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.DeleteEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.deleteEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.DeployKeyEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
		}.String())
		for _, hook := range bot.deployKeyEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.DeploymentEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.deploymentEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.DeploymentStatusEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.deploymentStatusEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.ForkEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.forkEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.GitHubAppAuthorizationEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
		}.String())
		for _, hook := range bot.gitHubAppAuthorizationEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.GollumEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.gollumEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.InstallationEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
		}.String())
		for _, hook := range bot.installationEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.InstallationRepositoriesEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
		}.String())
		for _, hook := range bot.installationRepositoriesEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.IssueCommentEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.issueCommentEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.IssueEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
		}.String())
		for _, hook := range bot.issueEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.IssuesEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.issuesEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.LabelEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
			Org:  *(e.GetOrg().Name),
			Repo: *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.labelEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.MarketplacePurchaseEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
		}.String())
		for _, hook := range bot.marketplacePurchaseEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.MemberEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.memberEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.MembershipEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
		}.String())
		for _, hook := range bot.membershipEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.MetaEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
		}.String())
		for _, hook := range bot.metaEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.MilestoneEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.milestoneEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.OrgBlockEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
		}.String())
		for _, hook := range bot.orgBlockEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.OrganizationEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
		}.String())
		for _, hook := range bot.organizationEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PageBuildEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.pageBuildEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PingEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
		}.String())
		for _, hook := range bot.pingEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.ProjectCardEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.projectCardEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.ProjectColumnEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.projectColumnEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.ProjectEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.projectEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PublicEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.publicEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PullRequestEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.pullRequestEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PullRequestReviewCommentEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.pullRequestReviewCommentEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PullRequestReviewEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.pullRequestReviewEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.PushEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.pushEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.ReleaseEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.releaseEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.RepositoryEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.repositoryEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.RepositoryVulnerabilityAlertEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
		}.String())
		for _, hook := range bot.repositoryVulnerabilityAlertEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.StarEvent:
		bot.logger.Println(eventTriggerLog{
			Type: typ,
		}.String())
		for _, hook := range bot.starEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.StatusEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.statusEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.TeamAddEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.teamAddEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.TeamEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Org:    *(e.GetOrg().Name),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.teamEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	case *github.WatchEvent:
		bot.logger.Println(eventTriggerLog{
			Type:   typ,
			Sender: *(e.GetSender().Login),
			Repo:   *(e.GetRepo().Name),
		}.String())
		for _, hook := range bot.watchEventHooks {
			if err := hook(ctx, e); err != nil {
				bot.logger.Printf("error on hook: %+v", err)
				return xerrors.Errorf("error on hook: %w", err)
			}
		}
	default:
		bot.logger.Printf("unsupported event type: %s", typ)
		return xerrors.Errorf("unsupported event type: %s", typ)
	}
	return nil
}

type eventTriggerLog struct {
	Type   string `json:"type"`
	Sender string `json:"sender,omitempty"`
	Org    string `json:"org,omitempty"`
	Repo   string `json:"repo,omitempty"`
}

func (etl eventTriggerLog) String() string {
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(etl); err != nil {
		panic(err)
	}
	return buf.String()
}

func (bot *Bot) AddCheckRunEventHook(hook CheckRunEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.checkRunEventHooks = append(bot.checkRunEventHooks, hook)
}

func (bot *Bot) AddCheckSuiteEventHook(hook CheckSuiteEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.checkSuiteEventHooks = append(bot.checkSuiteEventHooks, hook)
}

func (bot *Bot) AddCommitCommentEventHook(hook CommitCommentEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.commitCommentEventHooks = append(bot.commitCommentEventHooks, hook)
}

func (bot *Bot) AddCreateEventHook(hook CreateEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.createEventHooks = append(bot.createEventHooks, hook)
}

func (bot *Bot) AddDeleteEventHook(hook DeleteEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.deleteEventHooks = append(bot.deleteEventHooks, hook)
}

func (bot *Bot) AddDeployKeyEventHook(hook DeployKeyEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.deployKeyEventHooks = append(bot.deployKeyEventHooks, hook)
}

func (bot *Bot) AddDeploymentEventHook(hook DeploymentEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.deploymentEventHooks = append(bot.deploymentEventHooks, hook)
}

func (bot *Bot) AddDeploymentStatusEventHook(hook DeploymentStatusEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.deploymentStatusEventHooks = append(bot.deploymentStatusEventHooks, hook)
}

func (bot *Bot) AddForkEventHook(hook ForkEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.forkEventHooks = append(bot.forkEventHooks, hook)
}

func (bot *Bot) AddGitHubAppAuthorizationEventHook(hook GitHubAppAuthorizationEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.gitHubAppAuthorizationEventHooks = append(bot.gitHubAppAuthorizationEventHooks, hook)
}

func (bot *Bot) AddGollumEventHook(hook GollumEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.gollumEventHooks = append(bot.gollumEventHooks, hook)
}

func (bot *Bot) AddInstallationEventHook(hook InstallationEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.installationEventHooks = append(bot.installationEventHooks, hook)
}

func (bot *Bot) AddInstallationRepositoriesEventHook(hook InstallationRepositoriesEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.installationRepositoriesEventHooks = append(bot.installationRepositoriesEventHooks, hook)
}

func (bot *Bot) AddIssueCommentEventHook(hook IssueCommentEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.issueCommentEventHooks = append(bot.issueCommentEventHooks, hook)
}

func (bot *Bot) AddIssueEventHook(hook IssueEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.issueEventHooks = append(bot.issueEventHooks, hook)
}

func (bot *Bot) AddIssuesEventHook(hook IssuesEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.issuesEventHooks = append(bot.issuesEventHooks, hook)
}

func (bot *Bot) AddLabelEventHook(hook LabelEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.labelEventHooks = append(bot.labelEventHooks, hook)
}

func (bot *Bot) AddMarketplacePurchaseEventHook(hook MarketplacePurchaseEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.marketplacePurchaseEventHooks = append(bot.marketplacePurchaseEventHooks, hook)
}

func (bot *Bot) AddMemberEventHook(hook MemberEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.memberEventHooks = append(bot.memberEventHooks, hook)
}

func (bot *Bot) AddMembershipEventHook(hook MembershipEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.membershipEventHooks = append(bot.membershipEventHooks, hook)
}

func (bot *Bot) AddMetaEventHook(hook MetaEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.metaEventHooks = append(bot.metaEventHooks, hook)
}

func (bot *Bot) AddMilestoneEventHook(hook MilestoneEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.milestoneEventHooks = append(bot.milestoneEventHooks, hook)
}

func (bot *Bot) AddOrgBlockEventHook(hook OrgBlockEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.orgBlockEventHooks = append(bot.orgBlockEventHooks, hook)
}

func (bot *Bot) AddOrganizationEventHook(hook OrganizationEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.organizationEventHooks = append(bot.organizationEventHooks, hook)
}

func (bot *Bot) AddPageBuildEventHook(hook PageBuildEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.pageBuildEventHooks = append(bot.pageBuildEventHooks, hook)
}

func (bot *Bot) AddPingEventHook(hook PingEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.pingEventHooks = append(bot.pingEventHooks, hook)
}

func (bot *Bot) AddProjectCardEventHook(hook ProjectCardEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.projectCardEventHooks = append(bot.projectCardEventHooks, hook)
}

func (bot *Bot) AddProjectColumnEventHook(hook ProjectColumnEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.projectColumnEventHooks = append(bot.projectColumnEventHooks, hook)
}

func (bot *Bot) AddProjectEventHook(hook ProjectEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.projectEventHooks = append(bot.projectEventHooks, hook)
}

func (bot *Bot) AddPublicEventHook(hook PublicEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.publicEventHooks = append(bot.publicEventHooks, hook)
}

func (bot *Bot) AddPullRequestEventHook(hook PullRequestEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.pullRequestEventHooks = append(bot.pullRequestEventHooks, hook)
}

func (bot *Bot) AddPullRequestReviewCommentEventHook(hook PullRequestReviewCommentEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.pullRequestReviewCommentEventHooks = append(bot.pullRequestReviewCommentEventHooks, hook)
}

func (bot *Bot) AddPullRequestReviewEventHook(hook PullRequestReviewEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.pullRequestReviewEventHooks = append(bot.pullRequestReviewEventHooks, hook)
}

func (bot *Bot) AddPushEventHook(hook PushEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.pushEventHooks = append(bot.pushEventHooks, hook)
}

func (bot *Bot) AddReleaseEventHook(hook ReleaseEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.releaseEventHooks = append(bot.releaseEventHooks, hook)
}

func (bot *Bot) AddRepositoryEventHook(hook RepositoryEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.repositoryEventHooks = append(bot.repositoryEventHooks, hook)
}

func (bot *Bot) AddRepositoryVulnerabilityAlertEventHook(hook RepositoryVulnerabilityAlertEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.repositoryVulnerabilityAlertEventHooks = append(bot.repositoryVulnerabilityAlertEventHooks, hook)
}

func (bot *Bot) AddStarEventHook(hook StarEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.starEventHooks = append(bot.starEventHooks, hook)
}

func (bot *Bot) AddStatusEventHook(hook StatusEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.statusEventHooks = append(bot.statusEventHooks, hook)
}

func (bot *Bot) AddTeamAddEventHook(hook TeamAddEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.teamAddEventHooks = append(bot.teamAddEventHooks, hook)
}

func (bot *Bot) AddTeamEventHook(hook TeamEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.teamEventHooks = append(bot.teamEventHooks, hook)
}

func (bot *Bot) AddWatchEventHook(hook WatchEventHook) {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	bot.watchEventHooks = append(bot.watchEventHooks, hook)
}
