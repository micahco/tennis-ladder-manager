# Script

For my first user story, “As a league manager, I want to create player profiles, so that I can track their performance (ranking).”
Given a league manager is in the interactive terminal, when they enter the **player** command, then the user is prompted to enter the player’s unique username and the player is entered into the database. 

For my second user story, "As a league manager, I want enter match results, so that I can update the players’ rankings."
Given a league manager is in the interactive terminal, when they enter the **match** command, then the user is prompted to enter the two players usernames and then the results of each set.

For my third user story, "As a league manager, I want to view the ladder, so that I can see player rankings in order."
Given a league manager is in the interactive terminal interface when they enter the **ladder** command, then the application displays the current league rankings as an ordered list. 

For IH#1: When you start the interactive shell, you are greeted with a short message describing what the program does. Each command also has a description saying what that command does. For example: the match command is used to "Record match results" and the player command is used to "Manage players".

For IH#2: When the user enters the **player** command to add a new player to the league, the program gives instructions to the user and indicates that the user must assign each player a username, and that this user name is case-insensitive, which makes it easier for the user to later refer to that same player.

For IH#3: The user can provide an optional number argument to the **ladder** command to only show the first n rankings instead of the entire league.

For IH#4: Users familiar with command line interfaces will expect some commands such as *help*, *clear*, and *exit* to be available within the interactive terminal. Also, users expect to be able to use the keyboard shortcuts CTRL-C to exit and CTRL-L to clear the terminal.

For IH#5: This program does not include a full undo feature, but instead makes it easy to delete recently created entities that may have been created on accident. Likewise, if an entity is accidently removed, it is trivial to recreate.

For IH#6: Let's explore the path provided to the user for the *match* command:
    1. user enters the match command to enter a new match.
    2. user is prompted to enter the two players usernames. 
    3. user is prompted to enter the first sets scores for each player. 
    4. user is asked whether they wish to continue entering more sets. 
    5. after entering all sets scores, the winner is calculated and the match reuslts are displayed to the user.

For IH#7: Each command has optional arguments to expedite the process for power users. For example: a player can be created quickly by using the **player add username** command.

For IH#8: Commands that remove entries from the database will warn the user about the effects of what they are doing. Also, removal commands require users to confirm they want to continue with the operation, with the default option being *No* when the user presses enter.

The first quality attribute is **Simplicity**. My goal was for the application to have an intuitive user interface that the user can feel comfortable with. This is exemplified in the explicit path provided to create a match as described before for IH#6.

A second quality attribute is **Usability**. Commands have optional parameters. Having good defaults make it easier for users uncomfortable with command line tools to use the program. Whereas power users can get things done faster by using arguments with the commands.

A third quality attribute is **Correctness**. I want users to be able to trust the program to provide the same output every time. To help with this, I made it so usernames are case insensitive, so that if a user wants to capitalize the name of players, the program will still work as expected. This prevents the user from accidentally created two players with the same username, which could lead to unexpected 