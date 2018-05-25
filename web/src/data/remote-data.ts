export const IDLE = "checkout/remoteDataStatus/IDLE";
export const LOADING = "checkout/remoteDataStatus/LOADING";
export const FAILURE = "checkout/remoteDataStatus/FAILURE";
export const SUCCESS = "checkout/remoteDataStatus/SUCCESS";

type RemoteDataSuccess<Payload> = {
  status: typeof SUCCESS;
  data: Payload;
};

type RemoteDataFailure<Error> = {
  status: typeof FAILURE;
  error: Error;
};

type RemoteDataLoading = {
  status: typeof LOADING;
};

export type RemoteData<Error, Payload> =
  | { status: typeof IDLE }
  | RemoteDataLoading
  | RemoteDataFailure<Error>
  | RemoteDataSuccess<Payload>;

export const isSuccess = <Payload>(
  data?: RemoteData<any, Payload>
): data is RemoteDataSuccess<Payload> => !!data && data.status === SUCCESS;

export const isFailure = <Error, Payload>(
  data?: RemoteData<Error, Payload>
): data is RemoteDataFailure<Error> => !!data && data.status === FAILURE;

export const isLoading = <Error, Payload>(
  data?: RemoteData<Error, Payload>
): data is RemoteDataLoading => !!data && data.status === LOADING;

// This has to return a singleton, because we don't want to reload
// the component when this is returned from selector
const IDLE_DATA: RemoteData<any, any> = { status: IDLE };
const LOADING_DATA: RemoteData<any, any> = { status: LOADING };

export const idle = (): RemoteData<any, any> => IDLE_DATA;

export const load = (): RemoteData<any, any> => LOADING_DATA;

export const succeed = <Payload>(
  data: Payload
): RemoteDataSuccess<Payload> => ({ status: SUCCESS, data });

export const fail = <Error>(error: Error): RemoteData<Error, any> => ({
  status: FAILURE,
  error
});

export const unwrap = <Payload>(
  remoteData?: RemoteData<any, Payload>
): Payload | undefined => {
  if (isSuccess(remoteData)) {
    return remoteData.data;
  }
};

export const map = <A, B, Error>(
  fn: (a: A) => B | undefined,
  data: RemoteData<Error, A> | undefined
): RemoteData<Error, B> => {
  if (isSuccess(data)) {
    const b = fn(data.data);
    return b ? succeed(b) : IDLE_DATA;
  } else {
    return data || IDLE_DATA;
  }
};

export const recover = <A, Error>(
  fn: (err: Error) => A | undefined,
  data: RemoteData<Error, A> | undefined
): RemoteData<any, A> => {
  if (isFailure(data)) {
    const b = fn(data.error);
    return b ? succeed(b) : IDLE_DATA;
  } else {
    return data || IDLE_DATA;
  }
};

export const map2 = <A, B, C, Error>(
  fn: (a: A, b: B) => C | undefined,
  data1: RemoteData<Error, A> | undefined,
  data2: RemoteData<Error, B> | undefined
): RemoteData<Error, C> => {
  if (isSuccess(data1) && isSuccess(data2)) {
    const b = fn(data1.data, data2.data);
    return b ? succeed(b) : IDLE_DATA;
  }
  if (isFailure(data1)) {
    return data1;
  }
  if (isFailure(data2)) {
    return data2;
  }
  return IDLE_DATA;
};

export const flatMap = <A, B, Error>(
  fn: (a: A) => RemoteData<Error, B> | undefined,
  data: RemoteData<Error, A> | undefined
): RemoteData<Error, B> =>
  (isSuccess(data) ? fn(data.data) : data) || IDLE_DATA;

export const filterMap = <A, B, Error>(
  fn: (a: A) => RemoteData<Error, B> | undefined,
  list: A[]
): B[] => {
  return list.reduce(
    (res, item: A) => {
      const data = fn(item);
      if (isSuccess(data)) {
        res.push(data.data);
      }
      return res;
    },
    [] as B[]
  );
};

export const withDefault = <A, Error>(
  defaultValue: A,
  data: RemoteData<Error, A> | undefined
): A => (isSuccess(data) ? data.data : defaultValue);