/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

goog.exportSymbol('proto.payload.EmployeePayload', null, global);
goog.exportSymbol('proto.payload.EmployeePayload.Type', null, global);

/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.payload.EmployeePayload = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.payload.EmployeePayload, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.payload.EmployeePayload.displayName = 'proto.payload.EmployeePayload';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.payload.EmployeePayload.prototype.toObject = function(opt_includeInstance) {
  return proto.payload.EmployeePayload.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.payload.EmployeePayload} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.payload.EmployeePayload.toObject = function(includeInstance, msg) {
  var f, obj = {
    type: jspb.Message.getFieldWithDefault(msg, 1, 0),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
    jwt: jspb.Message.getFieldWithDefault(msg, 3, ""),
    shift: jspb.Message.getFieldWithDefault(msg, 4, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.payload.EmployeePayload}
 */
proto.payload.EmployeePayload.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.payload.EmployeePayload;
  return proto.payload.EmployeePayload.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.payload.EmployeePayload} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.payload.EmployeePayload}
 */
proto.payload.EmployeePayload.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.payload.EmployeePayload.Type} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setJwt(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setShift(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.payload.EmployeePayload.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.payload.EmployeePayload.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.payload.EmployeePayload} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.payload.EmployeePayload.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getJwt();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getShift();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.payload.EmployeePayload.Type = {
  UNKNOWN: 0,
  GET_OPEN_SHIFTS: 5,
  SHIFT_DETAILS: 1,
  PICK_UP_SHIFT: 2,
  CALL_OFF_SHIFT: 3,
  GET_MY_SCHEDULES: 4
};

/**
 * optional Type type = 1;
 * @return {!proto.payload.EmployeePayload.Type}
 */
proto.payload.EmployeePayload.prototype.getType = function() {
  return /** @type {!proto.payload.EmployeePayload.Type} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.payload.EmployeePayload.Type} value */
proto.payload.EmployeePayload.prototype.setType = function(value) {
  jspb.Message.setField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.payload.EmployeePayload.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.payload.EmployeePayload.prototype.setId = function(value) {
  jspb.Message.setField(this, 2, value);
};


/**
 * optional string jwt = 3;
 * @return {string}
 */
proto.payload.EmployeePayload.prototype.getJwt = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.payload.EmployeePayload.prototype.setJwt = function(value) {
  jspb.Message.setField(this, 3, value);
};


/**
 * optional string shift = 4;
 * @return {string}
 */
proto.payload.EmployeePayload.prototype.getShift = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.payload.EmployeePayload.prototype.setShift = function(value) {
  jspb.Message.setField(this, 4, value);
};


goog.object.extend(exports, proto.payload);
