"use strict";

angular.module('SearchApp', [])
	.controller('SearchController', function($scope, $http) {
		$scope.searchTerm
		$scope.suggestions = []
		$scope.articles = []

		$scope.search = function(term) {
			$http.get('/autocomplete?q=' + term).
			  success(function(data, status, headers, config) {
			  	console.log(data)
			  }).
			  error(function(data, status, headers, config) {
			    // called asynchronously if an error occurs
			    // or server returns response with an error status.
			  });
		}
	});