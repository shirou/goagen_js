import React from 'react';
import * as api from "./api_request.ts";
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';
import {createStore, combineReducers} from 'redux';
import {reducer as reduxFormReducer} from 'redux-form';
import { Values } from 'redux-form-website-template';

import { default as FieldLevelValidationForm } from './components/FieldLevelValidationForm.tsx';

const dest = document.getElementById('content');
const reducer = combineReducers({
  form: reduxFormReducer // mounted under "form"
});
const store = (createStore)(reducer);


const showResults = (values: any) => {
  api.UserCreate(values).then((result) =>{
    if (result.status !== 200){
      throw new Error('invalid paramater');
    }
    alert("OK!")
  }).catch((error) => {
    alert(error.message)
  });
};


ReactDOM.render(
  <Provider store={store}>
    <div>
      <h2>Form</h2>
      <FieldLevelValidationForm onSubmit={showResults} />
      <Values form="fieldLevelValidation" />
    </div>
  </Provider>,
  dest
);
