import React from 'react';
import * as api from "./api_request.js";
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';
import {createStore, combineReducers} from 'redux';
import {reducer as reduxFormReducer} from 'redux-form';
import { Values } from 'redux-form-website-template';

const dest = document.getElementById('content');
const reducer = combineReducers({
  form: reduxFormReducer // mounted under "form"
});
const store = (createStore)(reducer);


const showResults = (values) => {
  api.UserCreate(values).then((result) =>{
    if (result.status !== 200){
      throw new Error('invalid paramater');
    }
    return result.text();
  }).catch((error) => {
    alert(error.message)
  });
};

const FieldLevelValidationForm = require('./components/FieldLevelValidationForm').default;

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
