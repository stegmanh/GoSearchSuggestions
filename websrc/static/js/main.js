"use strict";

var demo = new Vue({
    el: '#app',
    data: {
        term: null,
        results: null,
        selected: null
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
        		self.results = JSON.parse(xhr.responseText)
	      }
	      xhr.send()
    	}
    }
})

Vue.filter('take', function(value, limit) {
	return value.slice(0, limit)
})

document.getElementById("search").addEventListener('keyup', function(e) {
	console.log(e)
})