# TODO for ejabberd API library

ejabberd command-line:

- Add prompt for password on oauth token generation with -P
- Refactor code to streamline the workflow / avoid duplication of processing between structure oauthToken and file.
  Oauthtoken should be probably serializable as is.
- Verbose mode to help debug request.
- Option to print value as text or JSon.
