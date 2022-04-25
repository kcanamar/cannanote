////////////////////////
// Setup - Import deps 
////////////////////////
const mongoose = require('mongoose');
const { Schema, model } = mongoose;

const entrySchema = new Schema(
    {
    username: String,
    strain: String,
    type: String,
    amount: String,
    consumption: String,
    description: String,
    date: { type: Date, default: Date.now},
    tags: [String],
    meta: {
        votes: Number,
        favs: Number,
        },
    }, 
    { timestamps: true }
);

const Entry = model("Entry", entrySchema);

//////////////////////
// Export
//////////////////////
module.exports = Entry;