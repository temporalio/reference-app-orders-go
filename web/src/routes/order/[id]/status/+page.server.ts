export const load = async (request) => {
    const { id } = request.params;
    const res = await fetch(`http://localhost:8234/api/v1/namespaces/default/workflows/${id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    const body = await res.json();
    return body;
}