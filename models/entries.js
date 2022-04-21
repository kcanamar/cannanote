////////////////////////
// Setup - Import deps 
////////////////////////
const mongoose = require('mongoose');
const { Schema, model } = mongoose;

const entrySchema = new Schema({
    strain: {
        type: String,
        required: true
    },
    amount: {
        type: String,
        required: true
    },
    method: {
        type: String,
        required: true
    },
    description: {
        type: String,
        required: true
    },
    tags: Array
});

const Entry = model("Entry", userSchema);

//////////////////////
// Export
//////////////////////
module.exports = Entry;