/* exported Script */
/* globals console, _, s, HTTP */

/** Global Helpers
 *
 * console - A normal console instance
 * _       - An underscore instance
 * s       - An underscore string instance
 * HTTP    - The Meteor HTTP object to do sync http calls
 */

class Script {
  /**
   * @params {object} request
   */
  prepare_outgoing_request({ request }) {
    // request.params            {object}
    // request.method            {string}
    // request.url               {string}
    // request.auth              {string}
    // request.headers           {object}
    // request.data.token        {string}
    // request.data.channel_id   {string}
    // request.data.channel_name {string}
    // request.data.timestamp    {date}
    // request.data.user_id      {string}
    // request.data.user_name    {string}
    // request.data.text         {string}
    // request.data.trigger_word {string}

    if (request.data.user_name === 'birthdaybot') {
        return;
    }

    // Prevent the request and return a new message
    if (request.data.text.match(/help$/) || request.data.text.match(/^@birthdaybot$/)) {
      return {
        message: {
          text: `Hi, I'm @birthdaybot and I keep track of everyone's birthday! :tada:\nTo get yourself added to the list, simply type *@birthdaybot born [YYYY-MM-DD]*, for example *@birthdaybot born 1971-09-10* in any public or private chat. It really is that easy :blush:`,
        }
      };
    }

    if (!request.data.text.match(/born .*$/)) {
        return;
    }

    request.url = request.url.replace('{Username}', request.data.user_name);
    request.headers['X-Api-Key'] = '';
    request.method = 'PUT';

    const birthday = request.data.text.match(/([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))/);
    if (!birthday || birthday.length < 2) {
        return {
          message: {
            text: `Sorry, I can't understand the format you provided. Please use \`born [YYYY-MM-DD]\`, for example \`born 1971-09-10\` :smile:`,
          }
        };
    }

    request.data = {
        birthday: birthday[1],
    };

    return request;
  }

  /**
   * @params {object} request, response
   */
  process_outgoing_response({ request, response }) {
    // request              {object} - the object returned by prepare_outgoing_request

    // response.error       {object}
    // response.status_code {integer}
    // response.content     {object}
    // response.content_raw {string/object}
    // response.headers     {object}

    // Return false will abort the response
    // Return empty will proceed with the default response process

    if (!response.content.ok) {
        return {
            content: {
                text: 'This is a bit embarrassing :flushed:... looks like something went wrong. Can you please tell @jan about this?',
                parseUrls: true,
                attachments: [{
                  "color": "#FF0000",
                  "title": "Oops",
                  "text": `Error Invoking API: ${response.content.message}`,
                }],
            },
        };
    }

    return {
        content: {
            text: 'Thanks for letting me know. You\'ve been added to the list! :blush:',
        }
    }
  }
}
