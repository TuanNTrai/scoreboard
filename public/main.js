
const CREATE_GAME = '/test_create';
const EDIT_GAME = '/test_edit';

// We test with this
const edit = {
	"id": "1mp",
	"basketCount": 3,
	"active": 42,
	"baskets": {
		"1": {
			"orderNum": 1,
			"par": 3,
			"scores": {
				"Miikka": {
					"score": 0,
					"ob": 0
				},
				"Player 2": {
					"score": 0,
					"ob": 0
				}
			}
		},
		"2": {
			"orderNum": 2,
			"par": 3,
			"scores": {
				"Miikka": {
					"score": 0,
					"ob": 0
				},
				"Player 2": {
					"score": 0,
					"ob": 0
				}
			}
		},
		"3": {
			"orderNum": 3,
			"par": 3,
			"scores": {
				"Miikka": {
					"score": 0,
					"ob": 0
				},
				"Player 2": {
					"score": 0,
					"ob": 0
				}
			}
		},
		"sasas": {
			"sadasd": {
				"name": 10
			}
		}
	}
};



let GAME = {};


function isUniq(name, arr) {
	code = true;
	arr.forEach(player => {
		if (name.toLowerCase() === player.name.toLowerCase()) {
			code = false;
			return code;
		}
	});
	return code;
}

function addPlayer(e) {
	e.preventDefault();
	if (this.player.length < 1) {
		this.errors.add = 'At least one character needed';
		return;
	}

	else if (this.player.length > 16) {
		this.errors.add = 'Max length is 16';
		return;
	}

	else if (!isUniq(this.player, this.selectedPlayers)) {
		this.errors.add = 'Player already exists';
		return;
	}

	this.selectedPlayers.push({ name: this.player, selected: true });
	this.player = '';
}

function toggleSelected(player) {
	player.selected = !player.selected
}

function deletePlayer(name) {
	this.selectedPlayers = this.selectedPlayers.filter(player => {
		console.log(player);
		return player.name != name;
	});
}

function start() {
	this.playersArr = [];
	this.selectedPlayers.forEach(player => {
		if (player.selected) {
			this.playersArr.push(player.name);
		}
	});

	if (!this.playersArr.length) {
		this.errors.start = "At least one player must be selected"
		return;
	}

	else if (this.playersArr.length > 5) {
		this.errors.start = "Max 5 players"
		return;
	}

	const query = {
		players: this.playersArr,
		basketCount: 3
	};

	postData(CREATE_GAME, query).then((data) => {
		console.log(data.status);
		console.log(data);
		this.course = data;
		localStorage.setItem('id', this.course.id);
		// window.location.pathname = 'games/' + this.course.id + '/' + this.course.active;
	});
}

function sendData() {
	let jee = {};
	console.log('REQUEST WITH', jee, this.course);

	postData(EDIT_GAME, this.course).then((data) => {
		console.log('FROM SERVER:', data);
		this.course = data;
	});
}

function join(e) {
	e.preventDefault();

	if (!this.gameID.length)
		return ;

	fetch('/games/' + this.gameID)
		.then((response) => {
			console.log(response.status);
			if (response.status != 200) {
				this.locked++;
				this.gameID = '';
				if (this.locked >= 3) {
					this.isDisabled = true;
				}
			}

			return response.json();
		})
		.then((data) => {
			this.course = data;
			console.log(data);
			localStorage.setItem('id', this.course.id);
		});
}


function incScore(player) {
	this.course.baskets[this.course.active].scores[player].score++;
	this.course.baskets[this.course.active].scores[player].total++;
}

function decScore(player) {
	if (this.course.baskets[this.course.active].scores[player].score > 1) {
		this.course.baskets[this.course.active].scores[player].score--;
		this.course.baskets[this.course.active].scores[player].total--;
	}
}

function deleteGame() {
	if (!confirm('Are you sure?'))
		return;
	localStorage.removeItem('id');
	localStorage.removeItem('active');
	this.course = {
		active: 0
	};
}

function prev() {

}

function next() {

}

// TODO: Hide from user
var app = new Vue({
	el: '#app',
	data: {
		errors: {
			start: '',
			add: ''
		},
		gameID: '',
		locked: 0,
		isDisabled: false,
		player: '',
		// TODO: Get from server
		selectedPlayers: [
			{ name: 'Miikka', selected: true },
			{ name: 'Sande', selected: true },
			{ name: 'Pasi', selected: true },
			{ name: 'Joni', selected: false },
		],
		playersArr: [],
		// Game object
		course: {}
	},
	methods: {
		addPlayer: addPlayer,
		toggleSelected: toggleSelected,
		deletePlayer: deletePlayer,
		start: start,
		sendData: sendData,
		incScore: incScore,
		decScore: decScore,
		deleteGame: deleteGame,
		join: join
	},
	computed: {
		selectedCount() {
			let count = 0;
			this.selectedPlayers.forEach(player => {
				if (player.selected) {
					count++;
				}
			});
			return count;
		},
	},
	created: function () {
		const id = localStorage.getItem('id');
		if (id == null)
			return;
		const URL = `/games/${id}`;
		if (id != null) {
			console.log('COOKIE', id);
			fetch(URL)
				.then((response) => {
					return response.json();
				})
				.then((data) => {
					this.course = data;
					// this.$forceUpdate();
					console.log(this.course);
				});
		} else {
			console.log('NO COOKIE');
		}
	}
});


async function postData(url = '', data = {}) {
	const response = await fetch(url, {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		credentials: 'same-origin',
		headers: {
			'Content-Type': 'application/json'
		},
		redirect: 'follow',
		referrerPolicy: 'no-referrer',
		body: JSON.stringify(data)
	});
	return response.json();
}
