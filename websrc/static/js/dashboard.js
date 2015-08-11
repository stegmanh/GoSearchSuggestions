"use strict";

var dashboard = new Vue({
    el: '#dashboard',
    data: {
		data: null
    },

    methods: {
        test: function() {
            var self = this
            console.log(self.data)
        }
    }
})