{% if include.header %}
{% assign header = include.header %}
{% else %}
{% assign header = "###" %}
{% endif %}
Run container in Kubernetes for specified project image from werf.yaml (build if needed)

{{ header }} Syntax

```shell
werf kube-run [options] [IMAGE_NAME] [-- COMMAND ARG...]
```

{{ header }} Examples

```shell
  # Run specified image
  $ werf kube-run --repo test/test application

  # Run interactive shell in the image
  $ werf kube-run --repo test/test -it -- sh

  # Run image with specified command
  $ werf kube-run --repo test/test -- /app/run.sh
```

{{ header }} Options

```shell
      --auto-pull-secret=true
            Automatically create docker config secret in the namespace and plug it via pod`s        
            imagePullSecrets for private registry access (default $WERF_AUTO_PULL_SECRET or true if 
            not specified)
      --cache-repo=[]
            Specify one or multiple cache repos with images that will be used as a cache. Cache     
            will be populated when pushing newly built images into the primary repo and when        
            pulling existing images from the primary repo. Cache repo will be used to pull images   
            and to get manifests before making requests to the primary repo.
            Also, can be specified with $WERF_CACHE_REPO_* (e.g. $WERF_CACHE_REPO_1=...,            
            $WERF_CACHE_REPO_2=...)
      --config=''
            Use custom configuration file (default $WERF_CONFIG or werf.yaml in working directory)
      --config-templates-dir=''
            Custom configuration templates directory (default $WERF_CONFIG_TEMPLATES_DIR or .werf   
            in working directory)
      --dev=false
            Enable development mode (default $WERF_DEV).
            The mode allows working with project files without doing redundant commits during       
            debugging and development
      --dev-branch='_werf-dev'
            Set dev git branch name (default $WERF_DEV_BRANCH or "_werf-dev")
      --dev-ignore=[]
            Add rules to ignore tracked and untracked changes in development mode (can specify      
            multiple).
            Also, can be specified with $WERF_DEV_IGNORE_* (e.g. $WERF_DEV_IGNORE_TESTS=*_test.go,  
            $WERF_DEV_IGNORE_DOCS=path/to/docs)
      --dir=''
            Use specified project directory where project’s werf.yaml and other configuration files 
            should reside (default $WERF_DIR or current working directory)
      --docker-config=''
            Specify docker config directory path. Default $WERF_DOCKER_CONFIG or $DOCKER_CONFIG or  
            ~/.docker (in the order of priority)
            Command needs granted permissions to read and pull images from the specified repo
      --dry-run=false
            Indicate what the command would do without actually doing that (default $WERF_DRY_RUN)
      --env=''
            Use specified environment (default $WERF_ENV)
      --extra-options=''
            Pass extra options to "kubectl run" command (default $WERF_EXTRA_OPTIONS)
      --final-repo=''
            Docker Repo to store only those stages which are going to be used by the Kubernetes     
            cluster, in other word final images (default $WERF_FINAL_REPO)
      --final-repo-container-registry=''
            Choose repo container registry for .
            The following container registries are supported: ecr, acr, default, dockerhub, gcr,    
            github, gitlab, harbor, quay.
            Default $WERF_FINAL_REPO_CONTAINER_REGISTRY or auto mode (detect container registry by  
            repo address).
      --final-repo-docker-hub-password=''
            Docker Hub password for  (default $WERF_FINAL_REPO_DOCKER_HUB_PASSWORD)
      --final-repo-docker-hub-token=''
            Docker Hub token for  (default $WERF_FINAL_REPO_DOCKER_HUB_TOKEN)
      --final-repo-docker-hub-username=''
            Docker Hub username for  (default $WERF_FINAL_REPO_DOCKER_HUB_USERNAME)
      --final-repo-github-token=''
            GitHub token for  (default $WERF_FINAL_REPO_GITHUB_TOKEN)
      --final-repo-harbor-password=''
            Harbor password for  (default $WERF_FINAL_REPO_HARBOR_PASSWORD)
      --final-repo-harbor-username=''
            Harbor username for  (default $WERF_FINAL_REPO_HARBOR_USERNAME)
      --final-repo-quay-token=''
            quay.io token for  (default $WERF_FINAL_REPO_QUAY_TOKEN)
      --follow=false
            Enable follow mode (default $WERF_FOLLOW).
            The mode allows restarting the command on a new commit.
            In development mode (--dev), werf restarts the command on any changes (including        
            untracked files) in the git repository worktree
      --git-work-tree=''
            Use specified git work tree dir (default $WERF_WORK_TREE or lookup for directory that   
            contains .git in the current or parent directories)
      --home-dir=''
            Use specified dir to store werf cache files and dirs (default $WERF_HOME or ~/.werf)
      --insecure-registry=false
            Use plain HTTP requests when accessing a registry (default $WERF_INSECURE_REGISTRY)
  -i, --interactive=false
            Enable interactive mode (default $WERF_INTERACTIVE or false if not specified)
      --kube-config=''
            Kubernetes config file path (default $WERF_KUBE_CONFIG, or $WERF_KUBECONFIG, or         
            $KUBECONFIG)
      --kube-config-base64=''
            Kubernetes config data as base64 string (default $WERF_KUBE_CONFIG_BASE64 or            
            $WERF_KUBECONFIG_BASE64 or $KUBECONFIG_BASE64)
      --kube-context=''
            Kubernetes config context (default $WERF_KUBE_CONTEXT)
      --log-color-mode='auto'
            Set log color mode.
            Supported on, off and auto (based on the stdout’s file descriptor referring to a        
            terminal) modes.
            Default $WERF_LOG_COLOR_MODE or auto mode.
      --log-debug=false
            Enable debug (default $WERF_LOG_DEBUG).
      --log-pretty=true
            Enable emojis, auto line wrapping and log process border (default $WERF_LOG_PRETTY or   
            true).
      --log-project-dir=false
            Print current project directory path (default $WERF_LOG_PROJECT_DIR)
      --log-quiet=false
            Disable explanatory output (default $WERF_LOG_QUIET).
      --log-terminal-width=-1
            Set log terminal width.
            Defaults to:
            * $WERF_LOG_TERMINAL_WIDTH
            * interactive terminal width or 140
      --log-verbose=false
            Enable verbose output (default $WERF_LOG_VERBOSE).
      --loose-giterminism=false
            Loose werf giterminism mode restrictions (NOTE: not all restrictions can be removed,    
            more info https://werf.io/documentation/advanced/giterminism.html, default              
            $WERF_LOOSE_GITERMINISM)
      --namespace=''
            Use specified Kubernetes namespace (default [[ project ]]-[[ env ]] template or         
            deploy.namespace custom template from werf.yaml or $WERF_NAMESPACE)
      --overrides=''
            Inline JSON to override/extend any fields in created Pod, e.g. to add imagePullSecrets  
            field (default $WERF_OVERRIDES)
      --platform=''
            Enable platform emulation when building images with werf. The only supported option for 
            now is linux/amd64.
      --pod=''
            Set created pod name (default $WERF_POD or autogenerated if not specified)
      --repo=''
            Docker Repo to store stages (default $WERF_REPO)
      --repo-container-registry=''
            Choose repo container registry.
            The following container registries are supported: ecr, acr, default, dockerhub, gcr,    
            github, gitlab, harbor, quay.
            Default $WERF_REPO_CONTAINER_REGISTRY or auto mode (detect container registry by repo   
            address).
      --repo-docker-hub-password=''
            Docker Hub password (default $WERF_REPO_DOCKER_HUB_PASSWORD)
      --repo-docker-hub-token=''
            Docker Hub token (default $WERF_REPO_DOCKER_HUB_TOKEN)
      --repo-docker-hub-username=''
            Docker Hub username (default $WERF_REPO_DOCKER_HUB_USERNAME)
      --repo-github-token=''
            GitHub token (default $WERF_REPO_GITHUB_TOKEN)
      --repo-harbor-password=''
            Harbor password (default $WERF_REPO_HARBOR_PASSWORD)
      --repo-harbor-username=''
            Harbor username (default $WERF_REPO_HARBOR_USERNAME)
      --repo-quay-token=''
            quay.io token (default $WERF_REPO_QUAY_TOKEN)
      --rm=true
            Remove pod and other created resources after command completion (default $WERF_RM or    
            true if not specified)
      --rm-with-namespace=false
            Remove also a namespace after command completion (default $WERF_RM_WITH_NAMESPACE or    
            false if not specified)
      --secondary-repo=[]
            Specify one or multiple secondary read-only repos with images that will be used as a    
            cache.
            Also, can be specified with $WERF_SECONDARY_REPO_* (e.g. $WERF_SECONDARY_REPO_1=...,    
            $WERF_SECONDARY_REPO_2=...)
  -Z, --skip-build=false
            Disable building of docker images, cached images in the repo should exist in the repo   
            if werf.yaml contains at least one image description (default $WERF_SKIP_BUILD)
      --skip-tls-verify-registry=false
            Skip TLS certificate validation when accessing a registry (default                      
            $WERF_SKIP_TLS_VERIFY_REGISTRY)
      --ssh-key=[]
            Use only specific ssh key(s).
            Can be specified with $WERF_SSH_KEY_* (e.g. $WERF_SSH_KEY_REPO=~/.ssh/repo_rsa,         
            $WERF_SSH_KEY_NODEJS=~/.ssh/nodejs_rsa).
            Defaults to $WERF_SSH_KEY_*, system ssh-agent or ~/.ssh/{id_rsa|id_dsa}, see            
            https://werf.io/documentation/reference/toolbox/ssh.html
  -S, --synchronization=''
            Address of synchronizer for multiple werf processes to work with a single repo.
            
            Default:
             - $WERF_SYNCHRONIZATION, or
             - :local if --repo is not specified, or
             - https://synchronization.werf.io if --repo has been specified.
            
            The same address should be specified for all werf processes that work with a single     
            repo. :local address allows execution of werf processes from a single host only
      --tmp-dir=''
            Use specified dir to store tmp files and dirs (default $WERF_TMP_DIR or system tmp dir)
  -t, --tty=false
            Allocate a TTY (default $WERF_TTY or false if not specified)
      --virtual-merge=false
            Enable virtual/ephemeral merge commit mode when building current application state      
            ($WERF_VIRTUAL_MERGE by default)
```
