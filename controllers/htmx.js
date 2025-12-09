////////////////////////
// Setup - Import deps
////////////////////////
const Entry = require('../models/entries');
///////////////////////
// Exports
///////////////////////

module.exports = {
    deleteConfirm
};
///////////////////////
// Declare Routes 
///////////////////////
// Delete Confirm Modal
async function deleteConfirm(req, res) {

    res.send(`
        <div id="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay" _="on click trigger closeModal"></div>
                <div class="modal-content">
                        <h1>You will be deleting your Note</h1>
                        <button class="btn danger" _="on click trigger closeModal">Cancel</button>
                </div>
        </div>
    `)
}
