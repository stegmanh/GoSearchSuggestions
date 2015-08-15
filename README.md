# GoSearchSuggestions

Hello! This is a project I started to learn Go. It's basically a port of an old assignment from C# to Go!

I hope to have a front end with query suggestions on every key press. For testing purposes I have commited the title.txt file which is must smaller than the wiki titles I will use in the completed assignment.

The program is seperated into two different parts. The crawler and the web server. THe crawler does the.. crawler and the web server can serve files.

The crawler works by crawling the CNN robotx.txt and starting from there. The messaging queue is build using redis and search and storage are built with postgres

The front end is build using the gorilla middlewhere and packages. Front end build with html/css and vue.js
