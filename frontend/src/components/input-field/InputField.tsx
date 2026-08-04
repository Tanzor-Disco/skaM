import './InputField.css'

interface InputFieldProps {
    fieldName: string
}
export default function InputField({ fieldName }: InputFieldProps) {
    return (
        <div className="input-field">
            <label htmlFor={fieldName}>{fieldName}</label>
            <input type="text" id={fieldName} name={fieldName.toLowerCase()} />
        </div>
    )
}
