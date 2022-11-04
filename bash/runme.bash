
: ${RUNME_DIR:=$HOME/.runme}
: ${RUNME_FILE:=.rm.md}

init() {
  declare desc="gets runme git repo"
  declare gitrepo=${1:-$RUNME_REPO}

  : ${gitrepo:? required}
  rm -rf ${RUNME_DIR} 2>/dev/null || true 
  # git clone ${RUNME_REPO} ${RUNME_DIR}
  git clone ${gitrepo} ${RUNME_DIR}
  choose
}

list() {
  declare desc="lists available block in markdown"
  _runme list --filename ${RUNME_FILE}
}

print() {
  declare desc="lists available block in markdown"
  _runme print --filename ${RUNME_FILE} "$@"
}

choose() {
  declare desc="choose an exercise"
  
  mdfile=$(cd ${RUNME_DIR}; ls -1 *.md|fzf --height=50% --layout=reverse )
  cp ${RUNME_DIR}/${mdfile} ${RUNME_FILE}
}

irun() {
  declare desc="interactively choose and run a block"
  block=$(_runme json --filename ${RUNME_FILE} \
    | _jq -r '.document[]|select(.name)|.name' \
    | _fzf --height=50% --layout=reverse --preview="runme print --filename ${RUNME_FILE} {}" \
    | xargs
  )
  _runme run --filename ${RUNME_FILE} ${block} "$@"
}
main() {
    [[ "$TRACE" ]] && set -x

    # default command is interactive run:
    if [[ $# -eq 0 ]]; then
      irun
    else
      cmd-export init
      cmd-export choose
      cmd-export list
      cmd-export print

      cmd-ns "" "$@"

      #  _runme "$@"
    fi
}