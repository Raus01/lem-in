# LEM-IN
a digital version of an ant farm

## Objective 

Create a program lem-in that will read from a file (describing the ants and the colony) given in the arguments.

Upon successfully finding the quickest path, lem-in will display the content of the file passed as argument and each move the ants make from room to room.

How does it work?

You make an ant farm with tunnels and rooms.
You place the ants on one side and look at how they find the exit.
You need to find the quickest way to get n ants across a colony (composed of rooms and tunnels).

At the beginning of the game, all the ants are in the room ##start. The goal is to bring them to the room ##end with as few moves as possible.
The shortest path is not necessarily the simplest.
Some colonies will have many rooms and many links, but no path between ##start and ##end.
Some will have rooms that link to themselves, sending your path-search spinning in circles. Some will have too many/too few ants, no ##start or ##end, duplicated rooms, links to unknown rooms, rooms with invalid coordinates and a variety of other invalid or poorly-formatted input. In those cases the program will return an error message ERROR: invalid data format. If you wish, you can elaborate a more specific error message (example: ERROR: invalid data format, invalid number of Ants or ERROR: invalid data format, no start room found).
You must display your results on the standard output in the following format :

number_of_ants
the_rooms
the_links

Lx-y Lz-w Lr-o ...
x, z, r represents the ants numbers (going from 1 to number_of_ants) and y, w, o represents the rooms names.

A room is defined by "name coord_x coord_y", and will usually look like "Room 1 2", "nameoftheroom 1 6", "4 6 7".

The links are defined by "name1-name2" and will usually look like "1-2", "2-5".

## Running the code

```go run main.go examplefilename.txt```