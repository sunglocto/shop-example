# shop-example
This is a very bad example of an online shop that you can make with [Gin](https://gin-gonic.com) and [Tailwind](https://tailwindcss.com)

This was mostly just a test of Tailwind. I'm not a very good designer, so excuse how terrible the UI looks. Maybe I'll make new stuff with this technology that I had no idea even existed.

To add a product, just follow the syntax of the other 2 products in the file, the attributes are self explanatory. The products will then show up on the root page. If you make a Stripe test account, create the products in Stripe then get payment links and put them in the Stripe attribute of your products. Once you click on the stripe button, it will redirect you to the relevant buy page.

Realistically if you build this up and put love and care into this, you could make a simple and easy to use webstore that even your grandmother can edit easily.

## Tailwind
Tailwind does not work well with Golang HTML templates, so I decided to whip up a couple of scripts that get all the classes used in all my templates, and then compiles the relevant CSS. I hate JavaScript, which is why I got the Tailwind CLI binary and did not touch the configuration or make a server. The relevant LSP did not work at all on Neovim for whatever reason, but it was pretty easy to find out what each class did and how to use them. Tailwind combined with Go's amazing template engine, and potentially using [htmx](https://htmx.org) for easy interactivity could make for a fantastic experience.


## Why
This approach also basically matches all of my web-dev frustrations. I hate making CSS, and Tailwind abstracts it so much I haven't even made any of my own CSS for this entire project. I hate JavaScript, so instead of using express or something for the backend I use Go, which is quickly becoming one of my favorite programming language. It's trivial for me to use the included standard library's JSON encoding and decoding features that allow me to make a struct, dump it into a JSON file or get a JSON file and grab the struct, and immediately use it in the project. This is an amateur process that you should **not** use in an actual project, please use a relational database. In an actual project, I would most likely use [Gorm](https://gorm.io) with a (No)SQL driver. I hate using frontend frameworks like React and Angular because I believe they overcomplicate the process. htmx is basically HTML on steroids and once again I don't need to touch a single line of JS.

Thank God.

## Issues
It's virtually impossible to add and remove products easily. You may have noticed that the ID of products is actually the index of the product in the Product slice, as when Go list files in a directory, they are in alphabetical order. This is a hack and this implementation can be obliterated by just having products with alphabetical or symbolic names, bringing about a host of issues when other people may edit products.

Because the ID of a product is its index in the Products slice, all products must be in order. This is extremely problematic, imagine this example:

Prod 0 -> Prod 1 -> Prod 2 -> Prod 3 -> Prod 4

Lets say something happened to product 3 and we delete it.

Prod 0 -> Prod 1 -> Prod 2 -> Prod 4

Now our items are not in order and the program will crash, because it's impossible for us to figure out this gap.
This would be eliminated if we used a `map` instead, with the key being a stringified ID and the value being a Product. We'd be able to iterate over it and easily add items to it.

Products don't have stock numbers and the review system has not been implemented but the bones are there.

It's good that we do not perform I/O on the directory every time a request is made, but we do so every minute. We can also perform a GET on /refresh/ (realistically this should be a POST) but the refresh operation has no rate limit or authorization. Meaning that a bad actor could simply spam this endpoint and potentially crash the server.

There is no graphical or easy way to add products and a user must know what to put in each attribute. The program does not add placeholder data so you can add a product that is fully empty:

`{}`

And the program will blindly accept it with no issues.

The program also shows every single product on the home page. If you were to add a high amount of products, the template engine could crash or be very slow due to having to make an enormous amount of iterations.

Speaking of, this program does not do any caching of requests. This can be solved by caching HTML responses with Redis, a volatile key value store database used to temporarily store commonly accessed data. Meaning if you had any actual customers accessing the site at once, it will crawl to its knees.

## Overall
This approach to creating websites is probably the best that I ca find for myself. It would be cool if I can make anything good with this stack. If you are wondering, the full stack is:

- Gin as backend framework (essential)
- Gorm as ORM
- HTMX as frontend framework (essential)
- Tailwind for CSS (essential)
- PostgresQL for database

If you want to have a particularly bad time, you can discard Gorm and Postgres to make a very simple website.
I'm not very fond of webdev regardless, but I need to get a job. So I'm learning, innit bruv.

I've been becoming very fond of Go, and it's probably going to become the compiled language of my choice (with Lua being the scripting language) of course everyone has different opinions, so don't grab a pitchfork just yet.
There are other Go libraries I'd love to test and use, such as [Bubble Tea](https://github.com/charmbracelet/bubbletea).
