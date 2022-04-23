////////////////////////
// Setup - Import deps
////////////////////////
const Entry = require('../models/entries');
const entrySeed = require('../models/seed')
///////////////////////
// Exports
///////////////////////

module.exports = {
    create,
    new: newEntry,
    show,
    index,
    edit,
    update,
    delete: deleteEntry,
    seed,
    like
};
///////////////////////
// Declare Routes 
///////////////////////
// todo find a way to sort mongodb
// Index
async function index(req, res) {
    try {
        let entryDocuments = await Entry.find({});
        res.render('index.ejs', {
            entries: entryDocuments
        });
    } catch(err) {
        res.send(err);
    }
};

// New
function newEntry(req, res) {
    res.render('new.ejs');
};

// Show
async function show(req, res) {
    try {
        let foundEntry = await Entry.findById(req.params.id);
        res.render('show.ejs', {
            entry: foundEntry
        });
    } catch(err) {
        res.send(err);
    }
};

// Edit
async function edit(req, res) {
    try {
        let entry = await Entry.findById(req.params.id);
        res.render('edit.ejs', {entry});
    } catch(err) {
        res.send(err);
    }
};

// Update
async function update(req, res) {
    try {
        await Entry.findByIdAndUpdate(req.params.id, req.body);
        res.redirect(`/cannanote/${req.params.id}`);
    } catch(err) {
        res.send(err);
    }
};

// Create
async function create(req, res) {
    try {
        let freshEntry = await Entry.create(req.body);
        freshEntry.save()
        res.redirect(`/cannanote/${freshEntry._id}`)
    } catch(err) {
        res.send(err);
    }
};

// Delete
async function deleteEntry(req, res) {
    try {
        await Entry.findByIdAndDelete(req.params.id);
        res.redirect('/cannanote')
    } catch(err) {
        res.send(err)
    }
}

// Seed
async function seed(req, res) {
    try {
        await Entry.deleteMany({});
        await Entry.create(entrySeed);
        res.redirect('/cannanote')
    } catch(err) {
        res.send(err)
    }
}

// todo work out functionality of like button
// Like
async function like(req, res) {
    try {
        let foundEntry = await Entry.findById(req.params.id);
        foundEntry.meta.favs += 1
        foundEntry.save()
        res.redirect(`/cannanote`)
    } catch(err) {
        res.send(err)
    }
}