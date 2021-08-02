const kitchenURL = 'http://localhost:8080/api/';

const kitchenService = {
    requestMenu: (): Promise<Response> => {
        return fetch(kitchenURL + 'returnMenu', {
            method: 'GET'
        });
    },

    requestFirebaseAccounts: (): Promise<Response> => {
        return fetch(kitchenURL + 'returnAccounts', {
            method: 'GET'
        });
    },

    prepareCourse: (data: any): Promise<Response> => {
        return fetch(kitchenURL + 'prepareCourse', {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }
};

export default kitchenService;
