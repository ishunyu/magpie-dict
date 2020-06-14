function search (){
	var xhttp;
	var searchTerm = document.getElementById('searchbox').value;
	xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function (){
		if (this.readyState == 4 && this.status == 200) {
			results = JSON.parse(this.responseText);
			createView(results.data);
		}
	};
	xhttp.open('GET', 'search?searchText=' + searchTerm, true);
	xhttp.send();
}

function searchKeyDown (ele){
	if (event.key === 'Enter') {
		search();
	}
}

function createView (data){
	const results = document.querySelector('#results');
	results.innerHTML = '';
	if (data.length === 0) {
		return alert('Sorry, we could not find anything matching the searched term');
	}
	for (let d of data) {
		const container = document.createElement('div');
		const metadata = createMetadataView(d);
		container.appendChild(metadata);
		const subs = document.createElement('div');
		const chineseLines = createChineseView(d);
		subs.appendChild(chineseLines);
		const englishLines = createEnglishView(d);
		subs.appendChild(englishLines);
		container.appendChild(subs);
		container.classList.add('returned-item');
		subs.classList.add('subs');
		results.appendChild(container);
	}
}

function createMetadataView (d){
	const episode = d.episode;
	const timestamp = d.subs.post.a.start;
	const metadataView = document.createElement('div');
	metadataView.innerHTML = `
    <span>Episode: ${episode}</span>
    <span>Timestamp: ${timestamp}</span>
`;
	metadataView.classList.add('meta-data');
	return metadataView;
}

function createChineseView (d){
	const chineseLines = document.createElement('div');
	const chineseTexts = [];
	for (let key in d.subs) {
		if (d.subs[key]['a'].text.length !== 0) {
			chineseTexts.push(d.subs[key]['a'].text);
		}
	}

	for (let line of chineseTexts) {
		const lineElement = document.createElement('p');
		lineElement.innerHTML = `- ${line}`;
		chineseLines.appendChild(lineElement);
	}
	chineseLines.classList.add('Chinese');
	return chineseLines;
}

function createEnglishView (d){
	const englishLines = document.createElement('div');
	const englishTexts = [];
	for (let key in d.subs) {
		if (d.subs[key]['b'].text.length !== 0) {
			englishTexts.push(d.subs[key]['b'].text);
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
