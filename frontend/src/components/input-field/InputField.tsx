import './InputField.css'

interface InputFieldProps {
    fieldName: string
	required: boolean
}
export default function InputField({ fieldName,required}: InputFieldProps) {
    return (
        <div className="input-field">
            <label htmlFor={fieldName}>{fieldName}</label>
            <input type="text" id={fieldName} name={fieldName.toLowerCase()} required={required} />
        </div>
    )
}
