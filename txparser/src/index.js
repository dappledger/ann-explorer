let utils = require('ethereumjs-util');
let abi = require('ethereumjs-abi');
let _ = require('underscore');

function isArray(type) {
  return type.lastIndexOf(']') === type.length - 1;
}

function getElementType(type) {
  const i = type.lastIndexOf('[');
  return type.substring(0, i);
}

function formatSingle(type, data) {
  let decodedData;
  if (isArray(type)) {
    // TODO: handle each array appropriately
    const elementType = getElementType(type);
    decodedData = _.map(data, function(data) {
      return formatSingle(elementType, data);
    });
  } else if (type.includes('bytes')) {
    const dataBuffer = Buffer.from(data, 'utf8');
    decodedData = dataBuffer.toString('hex');
  } else {
    decodedData = data.toString();
  }
  return decodedData;
}

/**
 * Decodes constructor args.
 *
 * @param {Object} contractABI - ABI of contract whose args to decode
 * @param {string} bytecode - Constructor args bytecode
 * @returns {Object} decodedArgs - Object representing decoded args with name, type, and data fields
 */
function decodeConstructorArgs(contractABI, bytecode) {
  const constructor = _.findWhere(contractABI, { 'type': 'constructor'});
  const inputNames = _.pluck(constructor.inputs, 'name');
  const inputTypes = _.pluck(constructor.inputs, 'type');
  let decoded = abi.rawDecode(inputTypes, new Buffer(bytecode, 'hex'));
  let decodedArgs = _.map(decoded, function(e, i) {
    const data = formatSingle(inputTypes[i], e);
    return { 'name': inputNames[i], 'type': inputTypes[i], 'data': data };
  });
  return decodedArgs;
}

/**
 * Generates constructor args bytecode based on input data.
 *
 * @param {Object[]} inputs - Array of objects with name, and type fields
 * @param {string} inputs[].name - Name of argument
 * @param {string} inputs[].type - Type of argument
 * @returns {string} bytecode - Constructor args bytecode
 */
function encodeConstructorArgs(inputs) {
  const inputTypes = _.pluck(inputs, 'type')
  const args = _.pluck(inputs, 'data')
  const encoded = abi.rawEncode(inputTypes, args);
  const bytecode = encoded.toString('hex');
  return bytecode;
}

function decodeFunctionArgs(contractABI, data) {
	const functions = _.where(contractABI, { 'type': 'function'});
	const callFunc = _.find(functions,function(f){
  		const selector = utils.bufferToHex(utils.sha3(f.name + "(" + _.pluck(f.inputs, 'type').join(",")+")")).substring(0,10);
  		return selector == data.substring(0,10);
	});
	if(!callFunc){
		return ;
	}
	const inputTypes = _.pluck(callFunc.inputs, 'type');
	const inputNames = _.pluck(callFunc.inputs, 'name');
	const args = data.substring(10);
	let decoded = abi.rawDecode(inputTypes, new Buffer(args, 'hex'));
	let decodedArgs = _.map(decoded, function(e, i) {
    	const data = formatSingle(inputTypes[i], e);
    	return { 'name': inputNames[i], 'type': inputTypes[i], 'data': data };
  	});

	return {
		callFunc: callFunc, 
		decodedArgs: decodedArgs
	};
}

function printTable(funcArgs){

	let args = funcArgs.decodedArgs;
	let header = "<thead><td>#</td><td>Name</td><td>Type</td><td>Data</td></thead>";
	let trs = args.map((a,i) => {
		return '<tr><td>'+(i+1)+'</td><td>'+a.name+'</td><td>'+a.type+'</td><td>'+a.data+'</td></tr>';
	});

	let table = '<table class="pure-table pure-table-horizontal">'+header+'<tbody>'+trs.join('')+'</tbody></table>';
	let output = '<div><label class="label label-primary label-sm">' + funcArgs.callFunc.name + '</label></div>' + table;
	return output;
}

module.exports = {
  decodeConstructorArgs: decodeConstructorArgs,
  encodeConstructorArgs: encodeConstructorArgs,
  decodeFunctionArgs: decodeFunctionArgs
};

window.decodeFunctionArgs = decodeFunctionArgs;
window.decodeFunctionArgsToTable = function(contractABI, data) {
	let args = decodeFunctionArgs(contractABI, data);
	return printTable(args);
}
