# Script

* Video shows and explains where and how a functional acceptance criterion is satisfied for a user story.
    * For my first user story, “As a league manager, I want to create player profiles, so that I can track their performance (ranking).”
    * Given a league manager is in the interactive terminal, when they enter the “player” command, then the user is prompted to enter the player’s unique username and the player is entered into the database. 

* Video shows and explains where and how a functional acceptance criterion is satisfied for a SECOND user story.
    * For my second user story, "As a league manager, I want enter match results, so that I can update the players’ rankings."
    * Given a league manager is in the interactive terminal, when they enter the "match” command, then the user is prompted to enter the players usernames and set results.

* Video shows and explains where and how a functional acceptance criterion is satisfied for a THIRD user story.
    * For my third user story, "As a league manager, I want to view the ladder, so that I can see player rankings in order."
    * Given a league manager is in the interactive terminal interface when they enter the "ladder” command, then the application displays the current league rankings in an ordered list. 

* Video shows and explains Inclusivity Heuristic #1.
    * For IH#1: When you start the interactive shell, you are greeted with a message describing who the program is for and how to use it. Each command describe why you would want to use it. For example: the *match* command tells the user that it is used to enter match results or remove the previous match.

* Video shows and explains Inclusivity Heuristic #2.
    * For IH#2: When the user enters the player command to add a new player to the league, the program gives instructions to the user and indicates that the user must assign each player a username, and that this user name is case-insensitive, which makes it easier for the user to later refer to that player.

* Video shows and explains Inclusivity Heuristic #3.
    * For IH#3: The user can provide an optional number argument to the *ladder* command to only show the first n rankings instead of the entire league.

* Video shows and explains Inclusivity Heuristic #4.
    * For IH#4: Users familiar with command line interfaces will expect a *help* command to show the various commands available to the user. Also, users expect to be able to use the keyboard shortcuts CTRL-C to exit and CTRL-L to clear the terminal.

* Video shows and explains Inclusivity Heuristic #5.
    * For IH#5: This program does not include a full undo feature, but instead makes it easy to delete recently created entities that may have been created on accident. Likewise, if an entity is accidently removed, it is trivial to recreate.

* Video shows and explains Inclusivity Heuristic #6.
    * For IH#6: Let's explore the path provided to the user for the *match* command:
    1. user enters the match command to enter a new match.
    2. user is prompted to enter the two players usernames. 
    3. user is prompted to enter the first sets scores for each player. 
    4. user is asked whether they wish to continue entering more sets. 
    5. after entering all sets scores, the winner is calculated and the match reuslts are displayed to the user.

* Video shows and explains Inclusivity Heuristic #7.
    * For IH#7: Each command has optional arguments to expedite the process for power users. For example: a player can be created quickly by using the *player add unique_username* command.

* Video shows and explains Inclusivity Heuristic #8.
    * For IH#8: Commands that remove entries from the database will warn the user about the effects of what they are doing. Also, removal commands require users to confirm they want to continue with the operation, with the default option being No when the user presses enter.

* Video shows and explains a quality attribute with a corresponding non-functional acceptance criterion.
    * First quality attribute is **Simplicity**. My goal was for the application to have an intuitive user interface that the user can feel comfortable with. This is exemplified in the explicit path provided to create a match as explained before.

* Video shows and explains a SECOND quality attribute.
    * A second quality attribute is **Usability**. Commands have optional parameters. Sensible defaults make it easier for users uncomfortable with command line tools to use the program. Whereas power users can get things done faster by using arguments with the commands.

* Video shows and explains a THIRD quality attribute.
    * A third quality attribute is **Correctness**. I want users to be able to trust the program to provide the same output every time. To help with this, I made it so usernames are case insensitive, so that if a user wants to capitalize the name of players, the program will still work as expected. This prevents the user from accidentally created two players with the same username, which could lead to unexpected 