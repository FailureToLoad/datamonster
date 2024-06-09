# datamonster

An app for managing the complexity of Kingdom Death: Monster campaigns.

## Requirements

Go 1.21.1  
nvm 0.39.7  
npm 21.4.0  

The site uses an alias to mitigate CORS issues. Make sure to add the following alias for localhost to /etc/hosts.

`127.0.0.1       localhost dev.local`

## Set up

Add the following variable to your `.bashrc`. I'm intending to secure this all via the container build process but for how threadbare the project is currently (and the fact that anyone reading this can't even get to the database) a test user and pass is fine.

`export CONN_STRING="postgres://app:Password1@DB CONTAINER IP ADDRESS GOES HERE:5432/records"`

To run the api, open a terminal and navigate to the api folder within this project. The command to execute it is `go run main.go`.

To run the site, open a terminal and navigate to the web folder within this project. First run `nvm use` to set the NPM version from the .nvmrc file. Launch the site with `npm run dev`
