
const CREATE_GAME = '/test_create';

// We test with this
const edit = {
	"id": "tigeri-jing-5",
	"basketCount": 3,
	"1": {
		"par": 3,
		"Tigerking": {
			"score": 6,
			"ob": 1
		},
		"Jing Jang": {
			"score": 1,
			"ob": 0
		}
	},
	"2": {
		"par": 4,
		"Tigerking": {
			"score": 7,
			"ob": 1
		},
		"Jing Jang": {
			"score": 1,
			"ob": 0
		}
	},
	"3": {
		"par": 5,
		"Tigerking": {
			"score": 8,
			"ob": 1
		},
		"Jing Jang": {
			"score": 1,
			"ob": 0
		}
	}
}

// With this we can create game
const startingData = {
	basketCount: 3,
	players: ['Tiger King', 'Ying Jang']
};

function toggleSelected(player) {
	player.selected = !player.selected
}

function start() {
	let playersArr = [];
	this.selectedPlayers.forEach(player => {
		if (player.selected) {
			playersArr.push(player.name);
		}
	});

	const query = {
		players: playersArr,
		basketCount: 3
	};

	postData(CREATE_GAME, query).then((data) => {
		console.log(data);
		this.course = data;
		this.active = 1;
		window.location.pathname = 'games/' + this.course.id + '/' + this.course.active;
	});
}

var app = new Vue({
	el: '#app',
	data: {
		active: 0,
		// TODO: Get from server
		selectedPlayers: [
			{name: 'Miikka', selected: true},
			{name: 'Sande', selected: false},
			{name: 'Pesukarhu', selected: true},
		],
		// Game object
		course: {}
	},
	methods: {
		toggleSelected: toggleSelected,
		start: start
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
