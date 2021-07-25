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
    }
};

export default kitchenService;
