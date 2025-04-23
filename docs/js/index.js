const TEST_DATA = [
    {
        timestamp: new Date(Date.UTC(2023, 9, 1, 12, 0, 0)).toISOString(),  // 2023-10-01T12:00:00.000Z
        actor: 'John Doe',
        service: 'Example Service',
        details: 'Initial changelog entry'
    },
    {
        timestamp: new Date(Date.UTC(2023, 9, 2, 14, 30, 0)).toISOString(), // 2023-10-02T14:30:00.000Z
        actor: 'Jane Smith',
        service: 'Example Service',
        details: 'Second changelog entry'
    },
    {
        timestamp: new Date(Date.UTC(2023, 9, 3, 9, 15, 0)).toISOString(),  // 2023-10-03T09:15:00.000Z
        actor: 'Alice Johnson',
        service: 'Another Service',
        details: 'Third changelog entry'
    },
];

const createRow = ({timestamp, actor, service, details}) => {
    const tr = document.createElement('tr');
    const timestampTd = document.createElement('td');
    const actorTd = document.createElement('td');
    const serviceTd = document.createElement('td');
    const detailsTd = document.createElement('td');

    timestampTd.textContent = timestamp;
    actorTd.textContent = actor;
    serviceTd.textContent = service;
    detailsTd.textContent = details;

    tr.appendChild(timestampTd);
    tr.appendChild(actorTd);
    tr.appendChild(serviceTd);
    tr.appendChild(detailsTd);

    return tr;
}

const fillTable = (tbody, data) => {
    if (!data) {
        return;
    }
    console.log('Filling table with data:', data);
    tbody.innerHTML = '';
    data.forEach(entry => {
        const row = createRow(entry);
        tbody.appendChild(row);
    });
}

const getData = () => {
    const url = 'http://localhost:8080/changes';
    return fetch(url)
        .then(response => response.json())
        .then(data => data.map(entry => ({
            timestamp: entry.Timestamp,
            actor: entry.Actor,
            service: entry.Service,
            details: entry.Details
        })))
        .catch(error => console.error('Error fetching data:', error));
}

const updateTable = (tbody) => {
    getData().then(data => fillTable(tbody, data))
        .catch(error => console.error('Error updating table:', error));
}

document.onreadystatechange = () => {
    if (document.readyState === 'complete') {
        const tbody = document.querySelector('tbody');
        updateTable(tbody);
        setInterval(() => updateTable(tbody), 5000); // Update every 5 seconds
    }
}