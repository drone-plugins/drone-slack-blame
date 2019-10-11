The following is an example template showing all the options available when
using the slack-blame plugin.

*Build*
  _Action_: {{ .Build.Action }}
  _Created_: {{ .Build.Created }}
  _DeployTo_: {{ .Build.DeployTo}}
  _Event_: {{ .Build.Event }}
  _FailedStages_:
  {{ with .Build.FailedStages }}
    {{ . }}
  {{ end }}
  _FailedSteps_:
  {{ with .Build.FailedSteps }}
    {{ . }}
  {{ end }}
  _Finished_: {{ .Build.Finished }}
  _Number_: {{ .Build.Number }}
  _Parent_: {{ .Build.Parent }}
  _PullRequest_: {{ .Build.PullRequest }}
  _Started_: {{ .Build.Started }}
  _Status_: {{ .Build.Status }}
  _SourceBranch_: {{ .Build.SourceBranch }}
  _Tag_: {{ .Build.Tag }}
  _TargetBranch_: {{ .Build.TargetBranch }}

*Repo*
  _DefaultBranch_: {{ .Repo.DefaultBranch }}
  _FullName_: {{ .Repo.FullName }}
  _Link_: {{ .Repo.Link }}
  _Name_: {{ .Repo.Name }}
  _Owner_: {{ .Repo.Owner }}
  _Private_: {{ .Repo.Private }}
  _RemoteURL_: {{ .Repo.RemoteURL }}
  _SCM_: {{ .Repo.SCM }}
  _Visibility_: {{ .Repo.Visibility }}

*Commit*
  _After_: {{ .Commit.After }}
  _Author_: {{ .Commit.Author }}
  _AuthorAvatar_: {{ .Commit.AuthorAvatar }}
  _AuthorEmail_: {{ .Commit.AuthorEmail }}
  _AuthorName_: {{ .Commit.AuthorName }}
  _Before_: {{ .Commit.Before }}
  _Branch_: {{ .Commit.Branch }}
  _Link_: {{ .Commit.Link }}
  _Message_: {{ .Commit.Message }}
  _Ref_: {{ .Commit.Ref }}
  _SHA_: {{ .Commit.SHA }}

*Stage*
  _Arch_: {{ .Stage.Arch }}
  _DependsOn_:
  {{ with .Stage.DependsOn }}
    {{ . }}
  {{ end }}
  _Finished_: {{ .Stage.Finished }}
  _Kind_: {{ .Stage.Kind }}
  _Machine_: {{ .Stage.Machine }}
  _Name_: {{ .Stage.Name }}
  _Number_: {{ .Stage.Number }}
  _OS_: {{ .Stage.OS }}
  _Started_: {{ .Stage.Started }}
  _Status_: {{ .Stage.Status }}
  _Type_: {{ .Stage.Type }}
  _Variant_: {{ .Stage.Variant }}
  _Version_: {{ .Stage.Version }}

*Step*
  _Name_: {{ .Step.Name }}
  _Number_: {{ .Step.Number }}

*SemVer*
  _Build_: {{ .SemVer.Build }}
  _Error_: {{ .SemVer.Error }}
  _Major_: {{ .SemVer.Major }}
  _Minor_: {{ .SemVer.Minor }}
  _Patch_: {{ .SemVer.Patch }}
  _Prerelease_: {{ .SemVer.Prerelease }}
  _Short_: {{ .SemVer.Short }}
  _Version_: {{ .SemVer.Version }}

*Slack*
