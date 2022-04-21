////////////////////////
// Setup - Import deps 
////////////////////////
const mongoose = require('mongoose')

// MongoDB & Mongoose 
mongoose.connect(process.env.DATABASE_URL, {
    useNewUrlParser: true,
    useUnifiedTopology: true
})
const db = mongoose.connection;
db.on('error', (err) => console.log(err.message + "yeah... that didn't work"))
db.on('connected', () => console.log('mongoose connected'))
db.on('disconnected', () => console.log('mongoose disconnected'))