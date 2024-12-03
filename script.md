My program provides a set of tools to help the user manage a tennis league.

It's main features are that it allows users to add players, enter match data, and view league rankings.

Microservice B validates match input.

Microservice C provides player stats.

Microservice D provides customizable artwork.

Microservice A was created by my teammate Josh Cantie and provides a random inspirational quote.

Each microservice is its own ZeroMQ server written in Python. My main program is written in Go and communicates with each microservice by sending and receiving messages.

Here on the left is my main program making a message request to the validator service. And in the right pane is the Python code for the validator microservice.

Show console logs for each microservice.
