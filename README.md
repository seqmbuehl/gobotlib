# Welcome to gobot

This game was created as a go + web tech teaching tool. The server is the first go code I've written, improvements are welcome and will gladly be pulled in.

# This repository is missing the 3 main files, necessary to make the program work. Checkout nathantsoi.github.io/slides/gobot.html for the accompanying code and talk

# Prereqs

 - homebrew, available at brew.sh or by running ```ruby -e "$(curl -fsSL https://raw.github.com/Homebrew/homebrew/go/install)"```

 - golang, available via homebrew by running ```brew install go```

# Deployment

Deploying this app to heroku is easy

 - sign up if you don't have an account https://id.heroku.com/signup

 - install the heroku toolbelt https://toolbelt.heroku.com/

 - setup heroku ```heroku login```

 - cd to the gobot directory if you're not already there and create a heroku app ```heroku create -b https://github.com/kr/heroku-buildpack-go.git```

 - enable websockets ```heroku labs:enable websockets```

 - deploy the app ```git push heroku master```

 - open the app in your browser with ```heroku open```

 - check the status of the app with ```heroku ps```

 - logs are available via ```heroku logs --tail```

# Development

 - run the app with ```PORT=3000 go run main.go```

 - when adding external dependencies, be sure to ```godep save```

# Legal

This header must be included, exactly as follows, in any redistributions or derivations:

2014, Nathan Tsoi

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 3 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
