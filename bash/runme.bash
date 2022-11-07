  
: ${RUNME_DIR:=$HOME/.runme}
: ${RUNME_FILE:=.rm.md}
: ${RUNME_REPO_PREFIX:=https://gist.github.com/lalyos/}

debug() {
  ((DEBUG)) && echo "[DEBUG] $@" 1>&2
}

init() {
  declare desc="gets runme git repo"
  declare gitrepo=${1:-$RUNME_REPO}

  : ${gitrepo:? required}
  # accept naked gistId
  [[ $gitrepo =~ "github" ]] || gitrepo=${RUNME_REPO_PREFIX%/}/${gitrepo}

  rm -rf ${RUNME_DIR} 2>/dev/null || true 
  # git clone ${gitrepo} ${RUNME_DIR}
  git clone ${gitrepo} ${RUNME_DIR} 2>/dev/null
  choose-exercise
}

list() {
  declare desc="lists available block in markdown"
  runme list --filename ${RUNME_FILE}
}

print() {
  declare desc="lists available block in markdown"
  block=$(choose-block)
  runme print --filename ${RUNME_FILE} ${block}
}

choose-exercise() {
  declare desc="choose an exercise"
  
  mdfile=$(
    cd ${RUNME_DIR}; ls -1 *.md \
    | fzf --prompt="Choose an exercise:" --height=50% --layout=reverse
  )
  cp ${RUNME_DIR}/${mdfile} ${RUNME_FILE}
}

choose-block() {
    # fzf runs in a subprocess, needs the faked fzf command/function
  export -f runme
  runme json --filename ${RUNME_FILE} \
    | jq -r '.document[]|select(.name)|.name' \
    | fzf  \
        --height=50% \
        --layout=reverse \
        --preview="runme print --filename ${RUNME_FILE} {}" 
}

command-list() {
  cmd=${1:-print}
  shift
  runme json --filename ${RUNME_FILE} \
    | jq -r '"'${cmd}' \(.document[]|select(.name)|.name) '"$@"'"'
}

inter() {
  export -f command-list runme
  export RUNME_FILE
  command=$( command-list run \
    | fzf  \
        --prompt="RUN> " \
        --header 'CTRL-r (run) CTRL-p (print)  CTR-d (dry-run)  CTR-c (exit)' \
        --bind 'ctrl-r:change-prompt(RUN> )+reload(command-list run)' \
        --bind 'ctrl-d:change-prompt(DRY> )+reload(command-list run --dry-run)' \
        --bind 'ctrl-p:change-prompt(PRINT> )+reload(command-list print)' \
        --bind 'ctrl-h:change-prompt(PRINT> )+reload(command-list print)' \
        --with-nth 2 \
        --height=50% \
        --layout=reverse \
        --preview="runme print --filename ${RUNME_FILE} {2}" \
        --preview-window 'right,40%' \
    )
    runme --filename ${RUNME_FILE} ${command}
}

irun() {
  declare desc="interactively choose and run a block"
  block=$(choose-block)
  runme run --filename ${RUNME_FILE} ${block} "$@"
}
shell() {
  declare desc="starts an intercative basher shell"

  # runs help in debug just to produce BASHENV file
  DEBUG=1 $SELF_EXECUTABLE ::: help >& /dev/null
  # uses BASHENV for interctive session
  bash --rcfile <( echo 'PS1="BASHER> "'; cat $(ls -1rt /tmp/bashenv.*|tail -1))
}

main() {
    [[ "$TRACE" ]] && set -x

    # default command is interactive run:
    if [[ $# -eq 0 ]]; then
      inter
    else
      cmd-export init
      cmd-export choose-exercise exercise
      cmd-export irun
      cmd-export list
      cmd-export print
      cmd-export shell
      cmd-export inter

      cmd-ns "" "$@"

      # runme "$@"
    fi
}