function search() {
	var xhttp;
	var searchTerm = document.getElementById('searchbox').value;
	var showID = getShowId();
	xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState == 4 && this.status == 200) {
			results = JSON.parse(this.responseText);
			createView(results.data);
		}
	};
	xhttp.open('GET', 'search?searchText=' + searchTerm + '&showID=' + showID, true);
	xhttp.send();
}

function searchKeyDown(ele) {
	if (event.key === 'Enter') {
		search();
	}
}

function getShowId() {
	var shows = document.getElementsByName('showlist');
	for (const show of shows) {
		if (show.checked) {
			return show.value;
		}
	}

	return '';
}

function createView(data) {
	const results = document.querySelector('#results');
	results.innerHTML = '';
	if (data.length === 0) {
		return alert('Sorry, we could not find anything matching the searched term');
	}
	for (let d of data) {
		const container = document.createElement('div');

		const metadata = createMetadataView(d);
		container.appendChild(metadata);

		const morePrevious = createMoreView(d, true);
		container.appendChild(morePrevious)
		
		const subs = document.createElement('div');
		const chineseLines = createChineseView(d);
		subs.appendChild(chineseLines);
		
		const englishLines = createEnglishView(d);
		subs.appendChild(englishLines);
		
		container.appendChild(subs);

		const moreNext = createMoreView(d, false);
		container.appendChild(moreNext)

		container.classList.add('returned-item');
		
		subs.classList.add('subs');
		results.appendChild(container);
	}
}

function createMetadataView(d) {
	const show = d.show
	const episode = d.episode;
	const timestamp = d.subs[0].a.start == '' ? d.subs[0].b.start : d.subs[0].a.start;
	const metadataView = document.createElement('div');
	metadataView.innerHTML = `
    <span>${show} | ${episode}</span>
    <span>Timestamp: ${timestamp}</span>
`;
	metadataView.classList.add('meta-data');
	return metadataView;
}

function createMoreView(d, previous) {
	const moreView = document.createElement('div');
	moreView.innerHTML = 'more'
	subid = previous ? d.subs[0].id : d.subs[d.subs.length - 1].id;
	moreView.setAttribute('subid', subid)
	moreView.setAttribute('expandtype', previous)
	moreView.onclick = function () {
		expand(moreView);
	}

	moreView.classList.add('more')
	return moreView
}

function createChineseView(d) {
	const chineseLines = document.createElement('div');
	const chineseTexts = [];
	for (let sub of d.subs) {
		if (sub['a'].text.length !== 0) {
			chineseTexts.push(sub['a'].text);
		}
	}

	for (let line of chineseTexts) {
		const lineElement = createText(line);
		chineseLines.appendChild(lineElement);
	}
	chineseLines.classList.add('Chinese');
	return chineseLines;
}

function createText(line) {
	const lineElement = document.createElement('p');
	lineElement.innerHTML = `- ${line}`;
	return lineElement;
}

function createEnglishView(d) {
	const englishLines = document.createElement('div');
	const englishTexts = [];
	for (let sub of d.subs) {
		if (sub['b'].text.length !== 0) {
			englishTexts.push(sub['b'].text);
		}
	}

	for (let line of englishTexts) {
		const lineElement = document.createElement('p');
		lineElement.innerHTML = `- ${line}`;
		englishLines.appendChild(lineElement);
	}
	englishLines.classList.add('English');
	return englishLines;
}

function expand(ele) {
	subId = ele.getAttribute('subid')
	type = ele.getAttribute('expandtype') == 'true'

	xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState == 4 && this.status == 200) {
			results = JSON.parse(this.responseText);
			updateSubs(ele, type, results)
		}
	};
	xhttp.open('GET', 'subs?id=' + subId + "&type=" + type, true);
	xhttp.send();
}

function updateSubs(ele, previous, results) {
	subs = results.subs;

	console.log("updateSubs subs:");
	console.log(ele)
	console.log(subs);
	
	if (subs.length < 1) {
		return;
	}
	
	chineseSubs = ele.parentElement.getElementsByClassName('Chinese')[0];
	englishSubs = ele.parentElement.getElementsByClassName('English')[0];

	if (previous) {
		for (i = subs.length - 1; i >= 0; i--) {
			sub = subs[i];

			if (sub['a'].text.length !== 0) {
				chineseSubs.prepend(createText(sub['a'].text))
			}

			if (sub['b'].text.length !== 0) {
				englishSubs.prepend(createText(sub['b'].text))
			}
		}

		ele.setAttribute('subid', subs[0].id)
	}
	else {
		for (let sub of subs) {
			if (sub['a'].text.length !== 0) {
				chineseSubs.appendChild(createText(sub['a'].text))
			}

			if (sub['b'].text.length !== 0) {
				englishSubs.appendChild(createText(sub['b'].text))
			}
		}

		ele.setAttribute('subid', subs[subs.length - 1].id)
	}
}

function updateShows() {
	var xhttp;
	xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState == 4 && this.status == 200) {
			results = JSON.parse(this.responseText);
			shows = results.shows
			
			message = ""
			shows.forEach(show => {
				message += show.name + ": " + show.episode + ", "
			});
			message = message.substring(0, message.length - 2)

			document.getElementById("show-data").innerHTML = message
		}
	};
	xhttp.open('GET', 'shows', true);
	xhttp.send();
}

window.onload = updateShows