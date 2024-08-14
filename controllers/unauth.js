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
    logout,
    newSignup,
    returnLogin,
}
///////////////////////
// Routes
///////////////////////

function main(req, res) {
    res.render("main.ejs")
}

async function signup(req, res) {
        const {username} = req.body
        const password = bcrypt.hashSync(req.body.password, 10)
        try {
            let newUser = await User.create({username, password})

            // login the new user
            login(req, res)

        } catch (error) {
            // handle non unique usernames
            res.render('signup.ejs', {notUnique: true})
        }
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
                // redirect to signup
                res.status(400).send({error: "User Doesn't Exist"})
            }
        }
    })
}

function logout(req, res) {
    req.session.destroy((err) => {
        res.redirect("/")
    })
}

function returnLogin(req, res) {
    res.render('login.ejs')
}

function newSignup(req, res) {
    res.render('signup.ejs', {notUnique: false})
}