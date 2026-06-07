# Lords or Larceny

## How to install and run
1. Copy this repo to your machine
`git clone https://github.com/KeshavBansal42/Lords_of_Larceny`

2. Navigate to the project folder
`cd Lords_of_Larceny`

3. Copy the .env.example to a .env file
`cp .env.example .env`

4. Run the docker compose
`docker compose up --build -d`

Now, your Lords of Larceny should be working, in case you are having an error from daemon about ports not being available make sure your localhost:5432 and :3000 are free