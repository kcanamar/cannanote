////////////////////////
// Setup - Import deps and create app object
////////////////////////
require('dotenv').config()
require('./config/db')
const express = require('express')
const middleware = require('./middleware/mid')
const CannaRouter = require('./routes/entires')
const indexRouter = require('./routes/index')
const app = express()
const PORT = process.env.PORT
//////////////////////
// Declare Middleware
//////////////////////
middleware(app)
///////////////////////
// Declare Routes and Routers 
///////////////////////
app.use("/", indexRouter)
app.use("/cannanote", CannaRouter)
///////////////////////////
// Server Listener
///////////////////////////
app.listen(PORT, () => console.log(`Who..Are..You?...${PORT}`))