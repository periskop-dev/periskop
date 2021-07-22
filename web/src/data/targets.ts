import { Dispatch } from "redux";
import * as RemoteData from "data/remote-data";
import { registerReducer } from "data/store"
import { Targets, TargetsState } from "data/types"

export const FETCH = "periskop/targets/FETCH"
export const FETCH_SUCCESS = "periskop/targets/FETCH_SUCCESS"
export const FETCH_FAILURE = "periskop/targets/FETCH_FAILURE"

export type TargetsAction =
  | { type: typeof FETCH; }
  | { type: typeof FETCH_SUCCESS; targets: Targets }
  | { type: typeof FETCH_FAILURE; error: any }

export const fetchTargets = () => {
  return (
    dispatch: Dispatch<TargetsAction>
  ) => {
    dispatch(fetchingTargets())

    return fetch(`${parseHostName()}targets/`).then(response => {
      return response
        .json()
        .then(targets => dispatch(fetchedTargetsSuccessfully(targets)))
        .catch(err => dispatch(fetchedTargetsFailed(err)))
    })
  }
}

export function fetchingTargets(): TargetsAction {
  return { type: FETCH }
}

export function fetchedTargetsSuccessfully(targets: Targets): TargetsAction {
  return { type: FETCH_SUCCESS, targets }
}

export function fetchedTargetsFailed(error: any): TargetsAction {
  return { type: FETCH_FAILURE, error }
}

const initialState: TargetsState = {
  targets: RemoteData.idle(),
}

function targetsReducer(state = initialState, action: TargetsAction) {
  switch (action.type) {
    case FETCH:
      return {
        ...state,
        targets: RemoteData.load()
      }
    case FETCH_SUCCESS:
      return {
        ...state,
        targets: RemoteData.succeed(action.targets)
      }
    case FETCH_FAILURE:
      return {
        targets: RemoteData.fail(action.error)
      }
    default:
      return state
  }
}

registerReducer("targetsReducer", targetsReducer)

const parseHostName = () => {
  let windowUrl = new URL(window.location.origin)
  if (windowUrl.hostname === "localhost") {
    windowUrl.port = "7777"
  }
  return windowUrl
}
