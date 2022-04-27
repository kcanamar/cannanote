////////////////////////
// Setup - Import deps and create app object
////////////////////////
require('dotenv').config()
require('./config/db')
const express = require('express')
const middleware = require('./middleware/mid')
const CannaRouter = require('./routes/entires')
const UnauthRouter = require('./routes/unauth')
const app = express()
const PORT = process.env.PORT || 3001
//////////////////////
// Declare Middleware
//////////////////////
middleware(app)
///////////////////////
// Declare Routes and Routers 
///////////////////////
app.use("/", UnauthRouter)
app.use("/cannanote", CannaRouter)
///////////////////////////
// Server Listener
///////////////////////////
app.listen(PORT, () => console.log(`Who..Are..You?...${PORT}`))