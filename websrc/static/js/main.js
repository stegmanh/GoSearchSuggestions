"use strict";

var demo = new Vue({
    el: '#app',
    data: {
        term: null,
        results: null
    },

    methods: {
    	autoComplete: function() {
			var self = this
			if (self.term.length < 1) {
				self.results = null
				return
			}
    		var xhr = new XMLHttpRequest()
	      xhr.open('GET', "/autocomplete?q=" + self.term)
	      xhr.onload = function () {
	      		console.log(JSON.parse(xhr.responseText))
        		self.results = JSON.parse(xhr.responseText)
	      }
	      xhr.send()
    	}
    }
})