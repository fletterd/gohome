var Constants = require('../constants.js');
var initialState = require('../initialState.js');

module.exports = function(state, action) {
    var newState = [];

    switch(action.type) {
    case Constants.ZONE_LOAD_ALL:
        break;

    case Constants.ZONE_LOAD_ALL_FAIL:
        break;

    case Constants.ZONE_LOAD_ALL_RAW:
        newState = action.data;
        break;

    default:
        newState = state || initialState().zones;
    }

    return newState;
};