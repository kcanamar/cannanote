////////////////////////
// Setup - Import deps
////////////////////////
const User = require("../models/user")
const bcrypt = require("bcryptjs")
///////////////////////
// Exports
///////////////////////
module.exports = {
    main,
    signup,
    login,
    logout
}
///////////////////////
// Routes
///////////////////////

function main(req, res) {
    res.render("main.ejs")
}

function signup(req, res) {
        const {username} = req.body
        const password = bcrypt.hashSync(req.body.password, 10)
        User.create({username, password}, (err, user) => {
            if (err) {
                res.status(400).send(err)
            } else {
                res.redirect("/")
            }
        })
}

function login(req, res) {
    const {username, password} = req.body
    User.findOne({username}, (err, user) => {
        if (err) {
            res.status(400).send(err)
        } else {
            if (user) {
                const pwcheck = bcrypt.compareSync(password, user.password)
                if (pwcheck) {
                    req.session.username = username
                    req.session.loggedIn = true
                    res.redirect("/cannanote")
                } else {
                    res.status(400).send({error: "Wrong Password"})
                }
            } else {
                res.status(400).send({error: "User Doesnt Exist"})
            }
        }
    })
}

function logout(req, res) {
    req.session.destroy((err) => {
        res.redirect("/")
    })
}