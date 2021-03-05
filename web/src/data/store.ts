import { Middleware, applyMiddleware, compose, createStore, combineReducers, ReducersMapObject, Reducer } from "redux";
import thunk from "redux-thunk";

const reducers: ReducersMapObject = {};

function configureStore() {
  const middlewares: Middleware[] = [];

  middlewares.push(thunk);

  const rootReducer = () => ({});

  // @ts-ignore
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
  return createStore(
    rootReducer,
    composeEnhancers(applyMiddleware(...middlewares))
  );
}

export function registerReducer(reducerName: string, reducer: Reducer<any>) {
  reducers[reducerName] = reducer;
  recreateReducers();
}

function recreateReducers() {
  const newReducer = combineReducers({
    ...reducers
  });
  store.replaceReducer(newReducer);
}

export const store = configureStore();