read_var() {
  if [ -z "$1" ]; then
    echo "environment variable name is required"
    return 1
  fi

  local ENV_FILE='.env'
  if [ -n "$2" ]; then
    ENV_FILE="$2"
  fi

  local VAR
  VAR=$(grep "$1" "$ENV_FILE" | xargs)
  IFS="=" read -ra VAR <<< "$VAR"
  echo "${VAR[1]}"
}