/* eslint-disable-next-line */
export const fetchItems = (params) => {
  console.log(params)
  const createPromise = response => new Promise((resolve, reject) => {
    return fetch('/api/clusters?page='+params.page+'&size='+params.limit)
      .then(response => {  
        if (response.ok) {
          return response;
        } else {
          reject(new Error('error'))
        }
      }, error => {
        reject(new Error(error.message))
      }).then(response => response.json())
      .then(json => {
        console.log(json)
        resolve(json)
      }).catch(function(error) {
        console.log(error);
    });
  })

  const { limit } = params;
  const offset = params.offset !== undefined ? params.offset : 0;
  const defaultResponse = {
    count: 0,
    results: [],
  };
  return createPromise(defaultResponse);
};
