#!/bin/bash
###
### STYLING UTILS
###
# colors
PRIMARY='\033[0;34m'
GOOD='\033[0;32m'
BAD='\033[0;31m'
NEUTRAL='\033[0;33m'
BOLD='\033[1m'
RESET='\033[0m'
function hide_cursor() {
    command printf '\e[?25l'
}
function show_cursor() {
    command printf '\e[?25h'
}
function echo() {
    command echo -e "$@"
}
spinner() {
    hide_cursor
    local pid=$!
    local delay=.1
    local spinstr='|/-\'
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf " %c  " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b"
    done
    printf "    \b\b\b\b"
    show_cursor
}
function confirm() {
    PROMPT="$1"
    while true; do
        read -p "$PROMPT" yn
        case $yn in
            [Yy]* ) return 0;;
            [Nn]* ) return 1;;
            * ) echo "$NEUTRAL (y/n).$RESET";;
        esac
    done
}
### end of styling

### banner (gpterm)
echo "
$PRIMARY
        GGGGGGGGGGGGGPPPPPPPPPPPPPPPPP   TTTTTTTTTTTTTTTTTTTTTTTEEEEEEEEEEEEEEEEEEEEEERRRRRRRRRRRRRRRRR   MMMMMMMM               MMMMMMMM
     GGG::::::::::::GP::::::::::::::::P  T:::::::::::::::::::::TE::::::::::::::::::::ER::::::::::::::::R  M:::::::M             M:::::::M
   GG:::::::::::::::GP::::::PPPPPP:::::P T:::::::::::::::::::::TE::::::::::::::::::::ER::::::RRRRRR:::::R M::::::::M           M::::::::M
  G:::::GGGGGGGG::::GPP:::::P     P:::::PT:::::TT:::::::TT:::::TEE::::::EEEEEEEEE::::ERR:::::R     R:::::RM:::::::::M         M:::::::::M
 G:::::G       GGGGGG  P::::P     P:::::PTTTTTT  T:::::T  TTTTTT  E:::::E       EEEEEE  R::::R     R:::::RM::::::::::M       M::::::::::M
G:::::G                P::::P     P:::::P        T:::::T          E:::::E               R::::R     R:::::RM:::::::::::M     M:::::::::::M
G:::::G                P::::PPPPPP:::::P         T:::::T          E::::::EEEEEEEEEE     R::::RRRRRR:::::R M:::::::M::::M   M::::M:::::::M
G:::::G    GGGGGGGGGG  P:::::::::::::PP          T:::::T          E:::::::::::::::E     R:::::::::::::RR  M::::::M M::::M M::::M M::::::M
G:::::G    G::::::::G  P::::PPPPPPPPP            T:::::T          E:::::::::::::::E     R::::RRRRRR:::::R M::::::M  M::::M::::M  M::::::M
G:::::G    GGGGG::::G  P::::P                    T:::::T          E::::::EEEEEEEEEE     R::::R     R:::::RM::::::M   M:::::::M   M::::::M
G:::::G        G::::G  P::::P                    T:::::T          E:::::E               R::::R     R:::::RM::::::M    M:::::M    M::::::M
 G:::::G       G::::G  P::::P                    T:::::T          E:::::E       EEEEEE  R::::R     R:::::RM::::::M     MMMMM     M::::::M
  G:::::GGGGGGGG::::GPP::::::PP                TT:::::::TT      EE::::::EEEEEEEE:::::ERR:::::R     R:::::RM::::::M               M::::::M
   GG:::::::::::::::GP::::::::P                T:::::::::T      E::::::::::::::::::::ER::::::R     R:::::RM::::::M               M::::::M
     GGG::::::GGG:::GP::::::::P                T:::::::::T      E::::::::::::::::::::ER::::::R     R:::::RM::::::M               M::::::M
        GGGGGG   GGGGPPPPPPPPPP                TTTTTTTTTTT      EEEEEEEEEEEEEEEEEEEEEERRRRRRRR     RRRRRRRMMMMMMMM               MMMMMMMM
$RESET"


echo "$NEUTRAL Checking if you have the required dependencies installed...$RESET"



function check_requirements() {
  echo "$NEUTRAL Checking for ollama...$RESET"
  sleep 1
  # check if ollama is installed
  if command -v ollama &> /dev/null
  then
      echo "ollama is installed"
  else
      echo "ollama is not installed"
      # prompt if user wants to install ollama
      if confirm "Would you like to install ollama? (y/n) "
      then
        echo "Installing ollama..."
        # install ollama
        PLATFORM=$(uname -s)
        LINUX="Linux"
        MAC="Darwin"
        WINDOWS="MINGW64_NT-10.0-19041"
        if [ "$PLATFORM" == "$LINUX" ]; then
          curl -s https://raw.githubusercontent.com/jmorganca/ollama/main/scripts/install.sh | bash
        elif [ "$PLATFORM" == "$MAC" ]; then
          ZIP="https://ollama.ai/download/Ollama-darwin.zip"
          curl -s $ZIP -o ollama.zip
          unzip ollama.zip
          rm ollama.zip
          mv ./Ollama.app /Applications
        elif [ "$PLATFORM" == "$WINDOWS" ]; then
          echo "Windows is not supported yet. Exiting..."
          exit 1
        else
          echo "Platform not supported. Exiting..."
          exit 1
        fi
      else
        echo "ollama is required to run gpterm. Exiting..."
        exit 1
      fi
  fi

  echo "$GOOD ollama is installed.$RESET"
  sleep 1
  echo "$NEUTRAL Checking for golang...$RESET"
  sleep 1

  # check if golang is installed
  if command -v go &> /dev/null
  then
      echo "go is installed"
  else
      echo "go is not installed"
      # prompt if user wants to install go
      if confirm "Would you like to install go? (y/n) "
      then
        echo "Installing go..."
        # install go
        PLATFORM=$(uname -s)
        LINUX="Linux"
        MAC="Darwin"
        WINDOWS="MINGW64_NT-10.0-19041"
        if [ "$PLATFORM" == "$LINUX" ]; then
          curl -s https://raw.githubusercontent.com/jmorganca/ollama/main/scripts/install.sh | bash
        elif [ "$PLATFORM" == "$MAC" ]; then
          ZIP="https://ollama.ai/download/Ollama-darwin.zip"
          curl -s $ZIP -o ollama.zip
          unzip ollama.zip
          rm ollama.zip
          mv ./Ollama.app /Applications
        elif [ "$PLATFORM" == "$WINDOWS" ]; then
          echo "Windows is not supported yet. Exiting..."
          exit 1
        else
          echo "Platform not supported. Exiting..."
          exit 1
        fi
      else
        echo "go is required to run gpterm. Exiting..."
        exit 1
      fi
  fi

  echo "$GOOD go is installed.$RESET"
  sleep 1
}

check_requirements & spinner

echo "$GOOD All dependencies are installed.$RESET"

echo "$NEUTRAL Installing gpterm...$RESET"

function install_gpterm() {
  GITURL="https://github.com/ProductionPanic/gpterm"
  sleep 1
  echo "$NEUTRAL Cloning gpterm from $GITURL...$RESET"
  git clone $GITURL gpterm
  cd gpterm
  echo "$NEUTRAL Building gpterm...$RESET"
  go build -o gpterm
  echo "$NEUTRAL Installing gpterm...$RESET"
  sudo mv gpterm /usr/local/bin
  echo "$GOOD gpterm is installed.$RESET"
  echo "$NEUTRAL Cleaning up...$RESET"
  cd ..
  rm -rf gpterm
}

install_gpterm & spinner

echo "$GOOD gpterm is installed.$RESET"

sleep 1

echo "$NEUTRAL Run gpterm by typing 'gpterm' in your terminal.$RESET"

GOODBYEBANNER="
$PRIMARY
    ::::::::  ::::::::: ::::::::::: :::::::::: :::::::::  ::::    ::::
   :+:    :+: :+:    :+:    :+:     :+:        :+:    :+: +:+:+: :+:+:+
   +:+        +:+    +:+    +:+     +:+        +:+    +:+ +:+ +:+:+ +:+
   :#:        +#++:++#+     +#+     +#++:++#   +#++:++#:  +#+  +:+  +#+
   +#+   +#+# +#+           +#+     +#+        +#+    +#+ +#+       +#+
   #+#    #+# #+#           #+#     #+#        #+#    #+# #+#       #+#
    ########  ###           ###     ########## ###    ### ###       ###

    Thanks for installing gpterm!
    Goodbye!
$RESET"

sleep 2

