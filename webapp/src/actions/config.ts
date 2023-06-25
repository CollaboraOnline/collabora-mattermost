import { AppDispatch } from '@/../types/store';
import Client from '@/client';
import Constants from '@/constants';

export const setConfig = (data: unknown) => {
  return async (dispatch: AppDispatch) => {
    dispatch({
      type: Constants.ACTION_TYPES.RECEIVED_CLIENT_CONFIG,
      data,
    });
  };
};

export const getConfig = () => {
  return async (dispatch: AppDispatch) => {
    let data = null;
    try {
      data = await Client.getConfig();
      dispatch(setConfig(data));
    } catch (error) {
      dispatch({
        type: Constants.ACTION_TYPES.CLIENT_CONFIG_ERROR,
        error,
      });

      return {
        data,
        error,
      };
    }

    return {
      data,
      error: null,
    };
  };
};
